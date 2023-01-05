package runner

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("shouldn't show anything when all outputs are disabled", func(t *testing.T) {
		server := httptest.NewServer(nil)
		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		req := New(server.Client(), nil, stdout, stderr, &Options{
			Method:       "post",
			Url:          server.URL,
			EnableColors: false,
		})

		req.Run()
		assert.Empty(t, stdout.String())
		assert.Empty(t, stderr.String())
	})

	t.Run("should show request headers", func(t *testing.T) {
		server := httptest.NewServer(nil)
		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		req := New(server.Client(), nil, stdout, stderr, &Options{
			Method:       "get",
			Url:          server.URL,
			ShowReqHead:  true,
			EnableColors: false,
		})

		req.Run()
		assert.Empty(t, stderr.String())

		outStr := stdout.String()
		assert.Contains(t, outStr, "GET / HTTP")
		assert.Contains(t, outStr, "Host: 127.0.0.1:")
		assert.Contains(t, outStr, "User-Agent: req")
		assert.Contains(t, outStr, "Accept: */*")
	})

	t.Run("should show response headers", func(t *testing.T) {
		server := httptest.NewServer(nil)
		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		req := New(server.Client(), nil, stdout, stderr, &Options{
			Method:       "get",
			Url:          server.URL,
			ShowRespHead: true,
			EnableColors: false,
		})

		req.Run()
		assert.Empty(t, stderr.String())

		outStr := stdout.String()
		assert.Contains(t, outStr, "Content-Type: text/plain")
	})

	t.Run("should use custom request header", func(t *testing.T) {
		server := httptest.NewServer(nil)
		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		req := New(server.Client(), nil, stdout, stderr, &Options{
			Method:       "get",
			Url:          server.URL,
			ShowReqHead:  true,
			EnableColors: false,
			Headers: map[string]string{
				"test":  "123",
				"test2": "123 456",
			},
		})

		req.Run()
		assert.Empty(t, stderr.String())

		outStr := stdout.String()
		assert.Contains(t, outStr, "GET / HTTP")
		assert.Contains(t, outStr, "Test: 123")
		assert.Contains(t, outStr, "Test2: 123 456")
	})

	t.Run("should show response body", func(t *testing.T) {
		json := `{ "id": 123, "name": "cool" }`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(json))
		}))

		stdout := &strings.Builder{}
		stderr := &strings.Builder{}

		req := New(server.Client(), nil, stdout, stderr, &Options{
			Method:       "put",
			Url:          server.URL,
			ShowRespBody: true,
			EnableColors: false,
		})

		req.Run()
		assert.Empty(t, stderr.String())

		outStr := stdout.String()
		assert.Contains(t, outStr, json)
	})

}

func TestParseOptions(t *testing.T) {
	t.Run("Wrong usage should not parse", func(t *testing.T) {
		_, err := ParseOptions([]string{})
		assert.Error(t, err)

		_, err = ParseOptions([]string{"./req"})
		assert.Error(t, err)

		_, err = ParseOptions([]string{"./req", "get"})
		assert.Error(t, err)
	})

	t.Run("Correct usage should parse properly", func(t *testing.T) {
		options, err := ParseOptions([]string{"./req", "get", "http://localhost:1234/"})
		assert.NoError(t, err)
		assert.Equal(t, "GET", options.Method)
		assert.Equal(t, "http://localhost:1234/", options.Url)
	})

	t.Run("should parse print options", func(t *testing.T) {
		options, err := ParseOptions([]string{"./req", "get", "http://localhost:1234/", "-p=Hb"})
		assert.NoError(t, err)
		assert.True(t, options.ShowReqHead)
		assert.False(t, options.ShowReqBody)
		assert.False(t, options.ShowRespHead)
		assert.True(t, options.ShowRespBody)
	})

	t.Run("should parse custom headers", func(t *testing.T) {
		options, err := ParseOptions([]string{"./req", "get", "http://localhost:1234/", "If-None-Match:123", "Another:hello guys"})
		assert.NoError(t, err)
		assert.Equal(t, "123", options.Headers["If-None-Match"])
		assert.Equal(t, "hello guys", options.Headers["Another"])
	})
}
