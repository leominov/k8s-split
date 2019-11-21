package main

import (
	"flag"
	"log"

	split "github.com/leominov/k8s-split"
)

var (
	input  = flag.String("f", "", "Path to file with Kubernetes specification")
	output = flag.String("o", "", "Path to output directory")
)

func main() {
	flag.Parse()
	err := split.Process(*input, *output)
	if err != nil {
		log.Fatal(err)
	}
}
