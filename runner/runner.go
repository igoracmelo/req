package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type ReqRunner struct {
	client  *http.Client
	logger  *log.Logger
	options *Options
}

func New(client *http.Client, logger *log.Logger, options *Options) *ReqRunner {
	return &ReqRunner{
		client:  client,
		logger:  logger,
		options: options,
	}
}

func (req *ReqRunner) Run() error {
	request, err := http.NewRequest(req.options.Method, req.options.Url, nil) // TODO: body
	if err != nil {
		return err
	}

	request.Header.Set("Host", request.Host)
	request.Header.Set("User-Agent", "req")
	request.Header.Set("Accept", "*/*")

	if req.options.ShowReqHead {
		req.logger.Printf("%s %s %s\n", request.Method, request.URL.Path, request.Proto)
		req.PrintHeaders(request.Header)
	}

	response, err := req.client.Do(request)
	if err != nil {
		return err
	}

	if req.options.ShowRespHead {
		req.logger.Println()
		req.logger.Printf("%s %s\n", response.Proto, response.Status)
		req.PrintHeaders(response.Header)
	}

	defer response.Body.Close()

	if req.options.ShowRespBody {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		req.logger.Println()
		req.logger.Println(string(body))
	}

	return nil
}

func (req *ReqRunner) PrintHeaders(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			req.logger.Printf("%s: %s\n", key, value)
		}
	}
}

type Options struct {
	Method       string
	Url          string
	ShowReqHead  bool
	ShowRespHead bool
	ShowRespBody bool
	EnableColors bool
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
		ShowRespHead: true,
		ShowRespBody: true,
		EnableColors: true,
	}

	return options, nil
}
