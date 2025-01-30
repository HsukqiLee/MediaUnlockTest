package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func AETV(c http.Client) Result {
	resp, err := GET(c, "https://geo.privacymanager.io/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		Country string `json:"country"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}
	if res.Country == "US" {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}
