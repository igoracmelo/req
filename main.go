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
var showReqHeaders = true
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

	req, err := http.NewRequest(options.Method, options.Url, options.Body)
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

	fmt.Println(color.Green(options.Method), color.Gray(req.URL.Path), color.Blue(req.Proto))

	if showReqHeaders {
		printHeaders(req.Header, &color)
		fmt.Println()
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(color.Cyan(resp.Proto), color.Yellow(resp.Status))

	if showRespHeaders {
		printHeaders(resp.Header, &color)
	}

	contentType := resp.Header.Get("Content-Type")

	if readBody {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body: ", err)
			os.Exit(1)
		}

		if showBody {
			fmt.Println()
			printRespBody(body, contentType, &color)
		}
	}

	fmt.Println()
}

func printRespBody(body []byte, contentType string, color *Color) error {
	if contentTypeIsCode(contentType) {
		err := quick.Highlight(os.Stdout, string(body), "go", "terminal256", "monokai")
		return err
	}

	isBinary := contentTypeIsBinary(contentType)

	if isBinary && showBinaryBody {
		fmt.Println(string(body))

	} else if isBinary {
		fmt.Println(color.Gray("[ Binary data ]"))

	} else {
		// plain text
		fmt.Println(string(body))
	}

	return nil
}

func printHeaders(headers http.Header, color *Color) {
	for name, values := range headers {
		for _, value := range values {
			// TODO: truncate on/off
			if len(value) > 80 {
				value = value[:80-3] + "..."
			}
			fmt.Printf("%s: %s\n", color.Blue(name), value)
		}
	}
}

type ReqOptions struct {
	Url     string
	Method  string
	Body    io.Reader
	Headers map[string]string
	Query   map[string]string
}

func parseOptions(args []string) (*ReqOptions, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("Not enough arguments")
	}

	options := &ReqOptions{
		Method:  strings.ToUpper(args[1]),
		Url:     args[2],
		Headers: map[string]string{},
		Query:   map[string]string{},
		Body:    strings.NewReader(""), // TODO:
	}

	options.Headers["Accept"] = "*/*"
	options.Headers["User-Agent"] = "req"

	for _, arg := range args[3:] {
		if strings.Contains(arg, ":") {
			chunks := strings.Split(arg, ":")
			options.Headers[chunks[0]] = chunks[1]

		} else if strings.Contains(arg, "==") {
			chunks := strings.Split(arg, "==")
			options.Query[chunks[0]] = chunks[1]
		}
	}

	return options, nil
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
