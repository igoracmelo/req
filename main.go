package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

type Color struct {
	Disable bool
}

func (c *Color) surround(text, color string) string {
	if c.Disable {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, "\033[0m")
}

func (c *Color) Red(text string) string {
	return c.surround(text, "\033[0;31m")
}

func (c *Color) Yellow(text string) string {
	return c.surround(text, "\033[0;33m")
}

func (c *Color) Blue(text string) string {
	return c.surround(text, "\033[0;34m")
}

func (c *Color) BBlue(text string) string {
	return c.surround(text, "\033[1;34m")
}

func (c *Color) Cyan(text string) string {
	return c.surround(text, "\033[0;36m")
}

func (c *Color) Gray(text string) string {
	return c.surround(text, "\033[0;37m")
}

var disableColors = false
var showReqHeaders = false // TODO:
var showRespHeaders = true
var readBody = true
var showBody = true
var showBinaryBody = false

func main() {
	// TODO: check
	method := strings.ToUpper(os.Args[1])
	url := os.Args[2]

	client := &http.Client{}

	color := Color{}

	if disableColors {
		color.Disable = true
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	for _, arg := range os.Args[3:] {
		// fmt.Println(arg)
		if strings.Contains(arg, ":") {
			chunks := strings.Split(arg, ":")
			req.Header.Set(chunks[0], chunks[1])
		} else if strings.Contains(arg, "==") {
			query := req.URL.Query()
			// req.URL.
			chunks := strings.Split(arg, "==")
			query.Set(chunks[0], chunks[1])
			req.URL.RawQuery = query.Encode()
		}
	}

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
			panic(err)
		}
	}

	fmt.Println()

	if showBody {
		bodyIsText := strings.HasPrefix(resp.Header.Get("Content-Type"), "text")

		if strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
			err := quick.Highlight(os.Stdout, string(body), "go", "terminal256", "monokai")
			if err != nil {
				panic(err)
			}
			return
		}

		if bodyIsText || showBinaryBody {
			fmt.Println(string(body))
		} else {
			fmt.Println(color.Gray("[ Binary data ]"))
		}
	}
}
