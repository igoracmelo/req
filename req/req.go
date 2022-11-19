package req

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Req struct {
	client  *http.Client
	logger  *log.Logger
	options *Options
}

func New(client *http.Client, logger *log.Logger, options *Options) *Req {
	return &Req{
		client:  client,
		logger:  logger,
		options: options,
	}
}

func (req *Req) Run() error {
	// client := http.Client{} // TODO: mock client

	request, err := http.NewRequest(req.options.Method, req.options.Url, nil) // TODO: body
	if err != nil {
		return err
	}

	request.Header.Set("Host", request.Host)
	request.Header.Set("User-Agent", "req")

	if req.options.ShowReqHeaders {
		req.PrintHeaders(request.Header)
	}

	response, err := req.client.Do(request)
	if err != nil {
		return err
	}

	if req.options.ShowRespHeaders {
		req.PrintHeaders(response.Header)
	}

	defer response.Body.Close()

	if req.options.ShowRespStatus {
		req.logger.Printf("%s %s\n", response.Proto, response.Status)
	}

	if req.options.ShowRespBody {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		req.logger.Println(string(body))
	}

	return nil
}

func (req *Req) PrintHeaders(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			req.logger.Printf("%s: %s\n", key, value)
		}
	}
}

// bytes.Buffer

type Options struct {
	Method          string
	Url             string
	ShowReqHeaders  bool
	ShowRespStatus  bool
	ShowRespHeaders bool
	ShowRespBody    bool
	EnableColors    bool
}

func ParseOptions(args []string) (*Options, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("Not enough args. Expected at least 3, got %d", len(args))
	}

	options := &Options{
		Method:          strings.ToUpper(args[1]),
		Url:             args[2],
		ShowReqHeaders:  true,
		ShowRespHeaders: true,
		ShowRespBody:    true,
	}
	return options, nil
}
