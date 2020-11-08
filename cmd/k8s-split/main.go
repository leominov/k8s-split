package main

import (
	"flag"
	"log"

	split "github.com/leominov/k8s-split"
)

var (
	input   = flag.String("f", "", "Path to file with Kubernetes specification")
	output  = flag.String("o", "", "Path to output directory")
	quiet   = flag.Bool("q", false, "Turn off k8s-split's output")
	splitby = flag.String("s", "", "Method to choose a directory's name: {prefix | tag}. Prefix will use longest object's name prefix whereas tag will use the value of app.kubernetes.io/part-of tag")
)

func main() {
	flag.Parse()
	split.Quiet = *quiet
	split.SplitBy = *splitby
	err := split.Process(*input, *output)
	if err != nil {
		log.Fatal(err)
	}
}
