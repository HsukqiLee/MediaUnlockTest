package mediaunlocktest

import (
	"net/http"
	"encoding/json"
	"io"
)

func NHKPlus(c http.Client) Result {
	resp, err := GET(c, "https://location-plus.nhk.jp/geoip/area.json")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
    
    b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
    var res struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	
	if res.CountryCode == "JP" {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusNo}
}