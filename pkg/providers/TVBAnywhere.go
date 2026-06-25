package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func TVBAnywhere(c http.Client) core.Result {
	resp, err := core.GET(c, "https://uapisfm.tvbanywhere.com.sg/geoip/check/platform/android")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res tvbAnywhereRes
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Country == "HK" || res.AllowInThisCountry {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.Country)}
	}
	return core.Result{Status: core.StatusNo}
}

type tvbAnywhereRes struct {
	AllowInThisCountry bool `json:"allow_in_this_country"`
	Country            string
}

