package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type LiveGetWebRequest struct {
}

func (receiver LiveGetWebRequest) FetchBytes(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "Conky-Currency-Rate")
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Wrong status code: %d", res.StatusCode)
	}
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func getRates(getWebRequest GetWebRequest, cfg CurrencySource) (ConvertResult, error) {
	var q string
	var result ConvertResult
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	for _, target := range cfg.Targets {
		q = q + "," + target + "_" + cfg.Parent
	}
	q = strings.Trim(q, ",")
	url := "https://free.currconv.com/api/v7/convert?q=" + q + "&compact=ultra&apiKey=" + cfg.Key
	body, err := getWebRequest.FetchBytes(ctx, url)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, err
}

func main() {
	var cfg CurrencySource
	flag.StringVar(&cfg.Parent, "parent", "RUB", "Code of parent currency")
	flag.StringVar(&cfg.Target, "targets", "USD,EUR", "Code of target currency, devided by comma")
	flag.Parse()
	cfg.Key = os.Getenv("CURRENCY_API")
	cfg.Targets = strings.Split(cfg.Target, ",")

	liveClient := LiveGetWebRequest{}
	result, err := getRates(liveClient, cfg)
	if err != nil {
		panic(err)
	}

	for s, f := range result {
		fmt.Printf("%s: %f", s, f)
		fmt.Println()
	}
}
