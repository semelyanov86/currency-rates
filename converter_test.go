package main

import (
	"testing"
)

type testWebRequest struct {
}

func (t testWebRequest) FetchBytes(url string) ([]byte, error) {
	return []byte(`{
  "USD_RUB": 46.1,
  "EUR_RUB": 33.4
}`), nil
}

func TestMainFunc(T *testing.T) {
	cfg := CurrencySource{
		Parent:  "RUB",
		Target:  "",
		Targets: []string{"EUR", "USD"},
		Key:     "12342134",
	}
	got, err := getRates(testWebRequest{}, cfg)
	if err != nil {
		T.Fatal(err)
	}
	if got["USD_RUB"] != 46.1 {
		T.Errorf("Expected contains %s, got %f", "USD_RUB: 46.1", got["USD_RUB"])
	}
	if got["EUR_RUB"] != 33.4 {
		T.Errorf("Expected contains %s, got %f", "EUR_RUB: 33.4", got["USD_RUB"])
	}
}
