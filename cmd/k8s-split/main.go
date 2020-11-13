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
	tag    = flag.Bool("tag", false, "Use the value of app.kubernetes.io/part-of tag as directory's name")
)

func main() {
	flag.Parse()
	if *tag && *prefix {
		log.Fatal("Choose either Prefix or Tag")
	}
	split.Quiet = *quiet
	split.Prefix = *prefix
	split.Tag = *tag
	err := split.Process(*input, *output)
	if err != nil {
		log.Fatal(err)
	}
}
