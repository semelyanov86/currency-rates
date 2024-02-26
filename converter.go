package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type LiveGetWebRequest struct {
}

type CBRDaily struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
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

func getRatesFromCbr(getWebRequest GetWebRequest, cfg CurrencySource) (ConvertResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	url := "https://www.cbr-xml-daily.ru/daily.xml"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyReader := strings.NewReader(string(bodyBytes))
	decoder := xml.NewDecoder(bodyReader)
	decoder.CharsetReader = charset.NewReaderLabel // This line sets the CharsetReader

	var cbrData CBRDaily
	if err := decoder.Decode(&cbrData); err != nil {
		return nil, err
	}

	// Assuming ConvertResult is a map[string]float64 where key is currency CharCode
	result := make(ConvertResult)
	for _, valute := range cbrData.Valutes {
		for _, target := range cfg.Targets {
			if valute.CharCode == target {
				// Convert value from string to float and adjust for nominal
				value, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", -1), 64)
				if err != nil {
					return nil, err
				}
				result[valute.CharCode] = value / float64(valute.Nominal)
			}
		}
	}

	return result, nil
}

func main() {
	var cfg CurrencySource
	flag.StringVar(&cfg.Parent, "parent", "RUB", "Code of parent currency")
	flag.StringVar(&cfg.Target, "targets", "USD,EUR", "Code of target currency, devided by comma")
	flag.Parse()
	cfg.Key = os.Getenv("CURRENCY_API")
	cfg.Targets = strings.Split(cfg.Target, ",")

	liveClient := LiveGetWebRequest{}
	result, err := getRatesFromCbr(liveClient, cfg)
	if err != nil {
		panic(err)
	}

	for s, f := range result {
		fmt.Printf("%s: %f", s, f)
		fmt.Println()
	}
}
