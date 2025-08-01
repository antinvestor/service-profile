package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pitabwire/frame"
)

type GeoIP struct {
	IP                 string  `json:"ip"`
	Network            string  `json:"network"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Country            string  `json:"country"`
	CountryName        string  `json:"country_name"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  int64   `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
}

func QueryIPGeo(ctx context.Context, svc *frame.Service, ip string) (*GeoIP, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	sts, resp, err := svc.InvokeRestService(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if sts != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d : %s", sts, string(resp))
	}

	var data GeoIP
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
