package split

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// Description Kubernetes specification
type Description struct {
	Kind     string
	Metadata struct {
		Name string
	}
}

// Process read inputFile, process and save to outputDir
func Process(inputFile, outputDir string) error {
	body, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}
	entries, err := ByEntries(body)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		kind, name, err := GetNameAndKind(entry)
		if err != nil {
			return err
		}
		log.Printf("Found %s.%s", name, kind)
		filename := path.Join(outputDir, fmt.Sprintf("%s.%s.yaml", name, kind))
		err = writeToFile(filename, entry)
		if err != nil {
			return err
		}
		log.Printf("Saved to %s", filename)
	}
	return nil
}

// ByEntries split multi-document YAML into separated maps
func ByEntries(body []byte) (result []map[string]interface{}, err error) {
	dec := yaml.NewDecoder(bytes.NewReader(body))
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
	out, err := yaml.Marshal(val)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, out, os.ModePerm)
	return err
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
