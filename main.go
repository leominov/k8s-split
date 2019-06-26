package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Description struct {
	Kind     string
	Metadata struct {
		Name string
	}
}

var (
	specsFile = flag.String("f", "", "Path to file with Kubernetes specification")
	outputDir = flag.String("o", "", "Path to output directory")
)

func main() {
	flag.Parse()
	body, err := ioutil.ReadFile(*specsFile)
	if err != nil {
		log.Fatal(err)
	}
	entries, err := SplitByEntries(body)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		kind, name, err := GetNameAndKind(entry)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found %s.%s", name, kind)
		filename := path.Join(*outputDir, fmt.Sprintf("%s.%s.yaml", name, kind))
		err = writeToFile(filename, entry)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Saved to %s", filename)
	}
}

func SplitByEntries(body []byte) (result []map[string]interface{}, err error) {
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
