package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func SBSonDemand(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.sbs.com.au/api/v3/network?context=odwebsite")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Get struct {
			Response struct {
				CountryCode string `json:"country_code"`
			} `json:"response"`
		} `json:"get"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if res.Get.Response.CountryCode == "AU" {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
