package main

import (
	"github.com/jarcoal/httpmock"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainFunc(T *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	key := os.Getenv("CURRENCY_API")

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://free.currconv.com/api/v7/convert?q=USD_RUB,EUR_RUB&compact=ultra&apiKey="+key,
		httpmock.NewStringResponder(200, `{
  "USD_RUB": 46.1,
  "EUR_RUB": 33.4
}`))
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if strings.Contains(string(out), "USD_RUB: 46.1") {
		T.Errorf("Expected contains %s, got %s", "USD_RUB: 46.1", string(out))
	}
	if strings.Contains(string(out), "EUR_RUB: 33.4") {
		T.Errorf("Expected contains %s, got %s", "EUR_RUB: 33.4", string(out))
	}
}
