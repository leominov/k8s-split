package main

import (
	"bytes"
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
	specsFile = flag.String("f", "", "Path to file with Kubernetes specifications")
	outputDir = flag.String("o", "", "Path to output directory")
)

func main() {
	flag.Parse()
	body, err := ioutil.ReadFile(*specsFile)
	if err != nil {
		log.Fatal(err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(body))
	for {
		var value map[string]interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		kind, name, err := getNameAndKind(value)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found %s.%s", name, kind)
		filename := path.Join(*outputDir, fmt.Sprintf("%s.%s.yaml", name, kind))
		err = writeToFile(filename, value)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Saved to %s", filename)
	}
}

func writeToFile(filename string, val interface{}) error {
	out, err := yaml.Marshal(val)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, out, os.ModePerm)
	return err
}

func getNameAndKind(val interface{}) (kind, name string, err error) {
	result := &Description{}
	if err := mapstructure.Decode(val, &result); err != nil {
		err = fmt.Errorf("Failed to decode body: %v", err)
	}
	kind = result.Kind
	name = result.Metadata.Name
	return
}
