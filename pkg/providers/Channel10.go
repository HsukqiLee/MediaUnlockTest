package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Channel10(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://10play.com.au/geo-web")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Allow bool `json:"allow"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		if strings.Contains(string(b), "not available") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Allow {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
