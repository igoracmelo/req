package req

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("shouldn't show anything when all outputs are disabled", func(t *testing.T) {
		server := httptest.NewServer(nil)
		out := &strings.Builder{}
		logger := log.New(out, "", 0)

		req := New(server.Client(), logger, &Options{
			Method: "GET",
			Url:    server.URL,
		})

		err := req.Run()
		assert.NoError(t, err)

		assert.Empty(t, out.String())
	})

	t.Run("should show request headers", func(t *testing.T) {
		server := httptest.NewServer(nil)
		out := &strings.Builder{}
		logger := log.New(out, "", 0)

		req := New(server.Client(), logger, &Options{
			Method:         "GET",
			Url:            server.URL,
			ShowReqHeaders: true,
		})

		err := req.Run()
		assert.NoError(t, err)

		assert.Contains(t, out.String(), "Host: 127.0.0.1:")
		assert.Contains(t, out.String(), "User-Agent: req")
	})

	t.Run("should show response headers", func(t *testing.T) {
		server := httptest.NewServer(nil)
		out := &strings.Builder{}
		logger := log.New(out, "", 0)

		req := New(server.Client(), logger, &Options{
			Method:          "GET",
			Url:             server.URL,
			ShowRespHeaders: true,
		})

		err := req.Run()
		assert.NoError(t, err)

		assert.Contains(t, out.String(), "Content-Type: text/plain")
	})
}

func TestParseOptions(t *testing.T) {
	t.Run("Wrong usage", func(t *testing.T) {
		_, err := ParseOptions([]string{})
		assert.Error(t, err)

		_, err = ParseOptions([]string{"./req"})
		assert.Error(t, err)

		_, err = ParseOptions([]string{"./req", "get"})
		assert.Error(t, err)
	})

	t.Run("Correct usage", func(t *testing.T) {
		options, err := ParseOptions([]string{"./req", "get", "http://localhost:1234/"})
		assert.NoError(t, err)
		assert.Equal(t, "GET", options.Method)
		assert.Equal(t, "http://localhost:1234/", options.Url)
	})
}
