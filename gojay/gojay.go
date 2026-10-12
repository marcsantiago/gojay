package main

import (
	"flag"
	"log"

	"github.com/marcsantiago/gojay/gojay/codegen"
)

//nolint:unused
var (
	pkg         = flag.String("pkg", "", "the package name of the generated file")
	dst         = flag.String("o", "", "destination file to output generated code")
	src         = flag.String("s", "", "source dir or file (absolute or relative path)")
	types       = flag.String("t", "", "types to generate")
	annotation  = flag.String("a", "json", "annotation tag (default json)")
	poolObjects = flag.String("p", "", "generate code to reuse objects using sync.Pool")
)

func main() {
	flag.Parse()
	options := codegen.NewOptionsWithFlagSet(flag.CommandLine)
	gen := codegen.NewGenerator(options)
	if err := gen.Generate(); err != nil {
		log.Fatal(err)
	}
}
