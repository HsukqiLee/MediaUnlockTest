package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Copilot(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://copilot.microsoft.com/c/api/user?api-version=2")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 302:
		if resp.Header.Get("Location") == "/" {
			return core.Result{Status: core.StatusBanned}
		}
		fallthrough
	case 403:
		return core.Result{Status: core.StatusNo}
	case 200:
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
	return core.Result{Status: core.StatusUnexpected}
}
