package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {

	t.Run("should err if doesn't have enough args", func(t *testing.T) {
		_, err := parseOptions([]string{"req", "get"})
		assert.Error(t, err)

		_, err = parseOptions([]string{"req", "google.com"})
		assert.Error(t, err)
	})

	t.Run("should parse method and url", func(t *testing.T) {
		args, err := parseOptions([]string{"req", "put", "google.com"})

		assert.NoError(t, err)
		assert.Equal(t, "google.com", args.Url)
		assert.Equal(t, "PUT", args.Method)
		assert.Empty(t, args.Headers)
		assert.Empty(t, args.Query)
	})

	t.Run("should parse request headers", func(t *testing.T) {
		args, err := parseOptions([]string{"req", "get", "site.com", "authorization:bearer 123", "x-test:cool"})

		assert.NoError(t, err)
		assert.Equal(t, "bearer 123", args.Headers["authorization"])
		assert.Equal(t, "cool", args.Headers["x-test"])
	})

	t.Run("should parse query string", func(t *testing.T) {
		args, err := parseOptions([]string{"req", "post", "nice.com", "sort==name", "limit==30"})

		assert.NoError(t, err)
		assert.Equal(t, "name", args.Query["sort"])
		assert.Equal(t, "30", args.Query["limit"])
	})
}
