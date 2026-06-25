package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func KKTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api.kktv.me/v3/ipcheck")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Data struct {
			Country   string
			IsAllowed bool `json:"is_allowed"`
		}
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Data.Country == "TW" && res.Data.IsAllowed {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}

