package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func Copilot(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://copilot.microsoft.com/c/api/user?api-version=2")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		var res struct {
			RegionCode string `json:"regionCode"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return core.Result{Status: core.StatusErr, Err: err}
		}
		if res.RegionCode != "" {
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.RegionCode)}
		}
	}
	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
