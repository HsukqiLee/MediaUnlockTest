package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Catchplay(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://sunapi.catchplay.com/geo",
		core.H{"authorization", "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy"},
		core.H{"accept", "application/json, text/plain, */*"},
		core.H{"accept-language", "zh-TW,zh;q=0.9,en;q=0.8"},
	)
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Code string `json:"code"`
		Data struct {
			IsoCode string `json:"isoCode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Code == "100016" {
		return core.Result{Status: core.StatusNo}
	}
	region := res.Data.IsoCode
	if region != "" {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
	}

	return core.Result{Status: core.StatusUnexpected}
}
