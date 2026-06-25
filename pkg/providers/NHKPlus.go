package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func NHKPlus(c http.Client) core.Result {
	resp, err := core.GET(c, "https://location-plus.nhk.jp/geoip/area.json")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	if res.CountryCode == "JP" {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusNo}
}

