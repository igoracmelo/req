package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/igoracmelo/req/color"
)

type ReqRunner struct {
	client  *http.Client
	logger  *log.Logger
	options *Options
	color   *color.Color
}

func New(client *http.Client, logger *log.Logger, options *Options) *ReqRunner {
	return &ReqRunner{
		client:  client,
		logger:  logger,
		options: options,
		color: &color.Color{
			Disable: !options.EnableColors,
		},
	}
}

func (r *ReqRunner) Run() error {
	r.options.Method = strings.ToUpper(r.options.Method)
	request, err := http.NewRequest(r.options.Method, r.options.Url, nil) // TODO: body
	if err != nil {
		return err
	}

	if request.URL.Path == "" {
		u, err := url.Parse(r.options.Url + "/")
		if err != nil {
			return err
		}
		request.URL = u
	}

	request.Header.Set("Host", request.Host)
	request.Header.Set("User-Agent", "req")
	request.Header.Set("Accept", "*/*")

	for key, value := range r.options.Headers {
		request.Header.Set(key, value)
	}

	if r.options.ShowReqHead {
		r.logger.Printf(
			"%s %s %s\n",
			r.color.Green(request.Method),
			r.color.Cyan(request.URL.Path),
			r.color.Blue(request.Proto),
		)
		r.PrintHeaders(request.Header)
	}

	response, err := r.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if r.options.ShowRespHead {
		r.logger.Println()
		r.logger.Printf("%s %s\n", response.Proto, r.color.Yellow(response.Status))
		r.PrintHeaders(response.Header)
	}

	if r.options.ShowRespBody {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		r.logger.Println()
		r.logger.Println(string(body))
	}

	return nil
}

func (r *ReqRunner) PrintHeaders(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			r.logger.Printf("%s: %s\n", r.color.Cyan(key), value)
		}
	}
}

type Options struct {
	Method       string
	Url          string
	ShowReqHead  bool
	ShowReqBody  bool
	ShowRespHead bool
	ShowRespBody bool
	EnableColors bool
	Headers      map[string]string
}

func ParseOptions(args []string) (*Options, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("Not enough args. Expected at least 3, got %d", len(args))
	}

	options := &Options{
		Method: strings.ToUpper(args[1]),
		Url:    args[2],

		// defaults
		ShowReqHead:  true,
		ShowReqBody:  true,
		ShowRespHead: true,
		ShowRespBody: true,
		EnableColors: true,
		Headers:      map[string]string{},
	}

	// TODO: refactor
	if len(args) > 3 {
		for _, arg := range args[3:] {
			if strings.HasPrefix(arg, "-") && strings.Contains(arg, "=") {
				chunks := strings.Split(arg[1:], "=")
				key := chunks[0]
				value := chunks[1]

				switch key {
				case "p":
					options.ShowReqHead = strings.Contains(value, "H")
					options.ShowReqBody = strings.Contains(value, "B")
					options.ShowRespHead = strings.Contains(value, "h")
					options.ShowRespBody = strings.Contains(value, "b")
				}
			} else if strings.Contains(arg, ":") {
				chunks := strings.Split(arg, ":")
				key := chunks[0]
				value := chunks[1]
				options.Headers[key] = value
			}
		}
	}

	return options, nil
}
