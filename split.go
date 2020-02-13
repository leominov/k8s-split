package split

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

var (
	Quiet bool
)

// List of Kubernetes specifications
type List struct {
	Kind  string
	Items []map[string]interface{}
}

// Description Kubernetes specification
type Description struct {
	Kind     string
	Metadata struct {
		Name string
	}
}

func readerFromInput(input string) (io.ReadSeeker, error) {
	if input == "-" {
		return os.Stdin, nil
	}
	r, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Process read input source, process and save to output directory
func Process(input, output string) error {
	r, err := readerFromInput(input)
	if err != nil {
		return err
	}
	entriesL, _ := ListByEntries(r)
	if len(entriesL) > 0 {
		return Save(entriesL, output)
	}
	r.Seek(0, 0)
	entriesM, err := MultiByEntries(r)
	if err != nil {
		return err
	}
	return Save(entriesM, output)
}

// Save save entries to output directory
func Save(entries []map[string]interface{}, output string) error {
	for _, entry := range entries {
		kind, name, err := GetNameAndKind(entry)
		if err != nil {
			return err
		}
		if !Quiet {
			log.Printf("Found %s.%s", name, kind)
		}
		filename := path.Join(output, fmt.Sprintf("%s.%s.yaml", name, kind))
		err = writeToFile(filename, entry)
		if err != nil {
			return err
		}
		if !Quiet {
			log.Printf("Saved to %s", filename)
		}
	}
	return nil
}

// ListByEntries split Kubernetes List into separated maps
func ListByEntries(r io.ReadSeeker) (result []map[string]interface{}, err error) {
	var l List
	err = yaml.NewDecoder(r).Decode(&l)
	if err != nil {
		return
	}
	result = append(result, l.Items...)
	return
}

// MultiByEntries split multi-document YAML into separated maps
func MultiByEntries(r io.ReadSeeker) (result []map[string]interface{}, err error) {
	dec := yaml.NewDecoder(r)
	for {
		var value map[string]interface{}
		err = dec.Decode(&value)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		result = append(result, value)
	}
	return
}

func writeToFile(filename string, val interface{}) error {
	out, _ := yaml.Marshal(val)
	return ioutil.WriteFile(filename, out, 0644)
}

// GetNameAndKind get Kubernetes `kind` and `name` from document
func GetNameAndKind(val interface{}) (kind, name string, err error) {
	result := &Description{}
	if err = mapstructure.Decode(val, &result); err != nil {
		err = fmt.Errorf("Failed to decode body: %v", err)
		return
	}
	kind = result.Kind
	if len(kind) == 0 {
		err = errors.New("Kind not found")
		return
	}
	name = result.Metadata.Name
	if len(name) == 0 {
		err = errors.New("Name not found")
		return
	}
	return
}
