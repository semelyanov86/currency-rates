package main

import "context"

type ConvertResult map[string]float64

type CurrencySource struct {
	Parent  string
	Target  string
	Targets []string
	Key     string
}

type GetWebRequest interface {
	FetchBytes(ctx context.Context, url string) ([]byte, error)
}
