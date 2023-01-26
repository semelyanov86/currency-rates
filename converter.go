package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type CurrencySource struct {
	Parent  string
	Target  string
	Targets []string
}

type ConvertResult map[string]float64

func main() {
	var cfg CurrencySource
	var convertResult ConvertResult
	flag.StringVar(&cfg.Parent, "parent", "RUB", "Code of parent currency")
	flag.StringVar(&cfg.Target, "targets", "USD,EUR", "Code of target currency, devided by comma")
	flag.Parse()
	key := os.Getenv("CURRENCY_API")
	cfg.Targets = strings.Split(cfg.Target, ",")
	var q string
	for _, target := range cfg.Targets {
		q = q + "," + target + "_" + cfg.Parent
	}
	q = strings.Trim(q, ",")
	result, err := http.Get("https://free.currconv.com/api/v7/convert?q=" + q + "&compact=ultra&apiKey=" + key)
	if err != nil {
		log.Fatal(err)
	}
	if result.StatusCode != http.StatusOK {
		log.Fatalf("Wrong status code: %d", result.StatusCode)
	}
	dec := json.NewDecoder(result.Body)
	if err := dec.Decode(&convertResult); err != nil {
		log.Fatalf("err - can not decode - %s", err)
	}
	for s, f := range convertResult {
		fmt.Printf("%s: %f", s, f)
		fmt.Println()
	}
}
