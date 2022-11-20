package main

import (
	"log"
	"net/http"
	"os"

	"github.com/igoracmelo/req/runner"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	run(os.Args, logger)
}

func run(args []string, logger *log.Logger) {
	options, err := runner.ParseOptions(args)

	if err != nil {
		showUsageAndExit(logger)
	}

	r := runner.New(&http.Client{}, logger, options)
	err = r.Run()
	_ = err // FIXME:
}

func showUsageAndExit(logger *log.Logger) {
	logger.Printf("usage: %s <get|post|etc> http://somesite.org/ [options]\n", os.Args[0])
	os.Exit(1)
}
