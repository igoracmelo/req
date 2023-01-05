package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/igoracmelo/req/runner"
)

func main() {
	run(os.Args, os.Stdin, os.Stdout, os.Stderr)
}

func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	options, err := runner.ParseOptions(args)

	if err != nil {
		fmt.Fprintf(stderr, "usage: %s <get|post|etc> http://somesite.org/ [options]\n", os.Args[0])
		os.Exit(1)
	}

	r := runner.New(&http.Client{}, stdin, stdout, stderr, options)
	err = r.Run()
	if err != nil {
		// TODO: error handling
		panic(err)
	}
}
