package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func HboGoAsia(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api2.hbogoasia.com/v1/geog?lang=undefined&version=0&bundleId=www.hbogoasia.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Country   string
		Territory string
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Territory == "" {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.Country)}
}

