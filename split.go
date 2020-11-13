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
	"strings"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

var (
	Quiet   bool
	Prefix bool
	Tag bool
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
		Name   string
		Labels struct {
			PartOf string `mapstructure:"app.kubernetes.io/part-of"`
		}
	}
}

func readerFromInput(input string) (io.ReadSeeker, error) {
	if input == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		return bytes.NewReader(b), err
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
	r.Seek(0, io.SeekStart)
	entriesM, err := MultiByEntries(r)
	if err != nil {
		return err
	}
	return Save(entriesM, output)
}

func findLongestNamePrefix(entries []map[string]interface{}) string {
	var names []string
	for _, entry := range entries {
		_, name, _, _ := GetNameAndKindAndPartof(entry)
		names = append(names, name)
	}
	return LongestCommonPrefix(names)
}

func preparePrefixedDirectory(entries []map[string]interface{}, output string) (string, error) {
	if Prefix {
		pref := findLongestNamePrefix(entries)
		if len(pref) == 0 {
			return output, nil
		}
		output = path.Join(output, pref)
		err := os.MkdirAll(output, 0755)
		if err != nil {
			return "", err
		}
	}
	if Tag {
		labels, err := FindUniqueLabelValues(entries)
		if err != nil {
			return "", err
		}
		for _, label := range labels {
			err := os.MkdirAll(path.Join(output, label), 0755)
			if err != nil {
				return "", err
			}
		}
	}

	return output, nil

}

// Save save entries to output directory
func Save(entries []map[string]interface{}, output string) error {
	prefixedDir, err := preparePrefixedDirectory(entries, output)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		kind, name, partof, err := GetNameAndKindAndPartof(entry)
		if err != nil {
			return err
		}
		if !Quiet {
			log.Printf("Found %s.%s", name, kind)
		}
		if Tag {
			output = path.Join(prefixedDir, partof)
		}
		if Prefix {
			output = prefixedDir
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
	if strings.ToLower(l.Kind) != "list" {
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

// GetNameAndKindAndPartof get Kubernetes `kind` and `name` from document
func GetNameAndKindAndPartof(val interface{}) (kind, name, partof string, err error) {
	result := &Description{}
	if err = mapstructure.Decode(val, &result); err != nil {
		err = fmt.Errorf("failed to decode body: %v", err)
		return
	}
	name = result.Metadata.Name
	if len(name) == 0 {
		err = errors.New("name not found")
		return
	}
	kind = result.Kind
	if len(kind) == 0 {
		err = errors.New("kind not found")
		return
	}
	partof = result.Metadata.Labels.PartOf

	return
}

// FindUniqueLabelValues returns list of unique app.kubernetes.io/part-of label in document
func FindUniqueLabelValues(entries []map[string]interface{}) ([]string, error) {
	var labels []string
	for _, entry := range entries {
		if _, _, label, err := GetNameAndKindAndPartof(entry); err == nil {
			labels = append(labels, label)
		}
	}

	j := 0
	for i := 1; i < len(labels); i++ {
		if labels[j] == labels[i] {
			continue
		}
		j++
		labels[j] = labels[i]
	}
	result := labels[:j+1]

	return result, nil
}
