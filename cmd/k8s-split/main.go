package main

import (
	"flag"
	"log"

	split "github.com/leominov/k8s-split"
)

var (
	input  = flag.String("f", "", "Path to file with Kubernetes specification")
	output = flag.String("o", "", "Path to output directory")
	quiet  = flag.Bool("q", false, "Turn off k8s-split's output")
	prefix = flag.Bool("prefix", false, "Use longest name prefix as directory's name")
)

func main() {
	flag.Parse()
	split.Quiet = *quiet
	split.Prefix = *prefix
	err := split.Process(*input, *output)
	if err != nil {
		log.Fatal(err)
	}
}
