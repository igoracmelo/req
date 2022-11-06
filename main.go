package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

var disableColors = false
var showReqHeaders = false // TODO:
var showRespHeaders = true
var readBody = true
var showBody = true
var showBinaryBody = false

func main() {
	options, err := parseOptions(os.Args)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	color := Color{}

	if disableColors {
		color.Disable = true
	}

	req, err := http.NewRequest(options.Method, options.Url, nil) // TODO: body
	if err != nil {
		panic(err)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	query := req.URL.Query()

	for key, value := range options.Query {
		query.Set(key, value)
	}

	req.URL.RawQuery = query.Encode()

	// fmt.Println(req.URL.String())
	// fmt.Println()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(color.Cyan(resp.Proto), color.Yellow(resp.Status))

	if showRespHeaders {
		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", color.Blue(name), value)
			}
		}
	}

	body := []byte{}

	if readBody {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body: ", err)
			os.Exit(1)
		}
	}

	fmt.Println()

	if showBody {
		contentType := resp.Header.Get("Content-Type")

		if contentTypeIsCode(contentType) {
			err := quick.Highlight(os.Stdout, string(body), "go", "terminal256", "monokai")
			if err != nil {
				panic(err)
			}
		} else if showBinaryBody || !contentTypeIsBinary(contentType) {
			fmt.Println(string(body))
		} else {
			fmt.Println(color.Gray("[ Binary data ]"))
		}
	}

	fmt.Println()
}

type ReqOptions struct {
	Url     string
	Method  string
	Body    string
	Headers map[string]string
	Query   map[string]string
}

func parseOptions(args []string) (*ReqOptions, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("Not enough arguments")
	}

	reqArgs := &ReqOptions{
		Method:  strings.ToUpper(args[1]),
		Url:     args[2],
		Headers: map[string]string{},
		Query:   map[string]string{},
	}

	for _, arg := range args {
		if strings.Contains(arg, ":") {
			chunks := strings.Split(arg, ":")
			reqArgs.Headers[chunks[0]] = chunks[1]

		} else if strings.Contains(arg, "==") {
			chunks := strings.Split(arg, "==")
			reqArgs.Query[chunks[0]] = chunks[1]
		}
	}

	return reqArgs, nil
}

func contentTypeIsCode(contentType string) bool {
	prefixes := []string{
		"text/html",
		"text/xml",
		"application/json",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(contentType, prefix) {
			return true
		}
	}

	return false
}

func contentTypeIsBinary(contentType string) bool {
	prefixes := []string{
		"text/",
		"application/json",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(contentType, prefix) {
			return false
		}
	}

	return true
}
