package mediaunlocktest

import (
    "io"
	"net/http"
	"encoding/json"
)

func SBSonDemand(c http.Client) Result {
	resp, err := GET(c, "https://www.sbs.com.au/api/v3/network?context=odwebsite")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		Get struct {
			Response struct {
			    CountryCode string `json:"country_code"`
			} `json:"response"`
		} `json:"get"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	if res.Get.Response.CountryCode == "AU" {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}