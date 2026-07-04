package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Coze(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.coze.com/api/developer/get_login_info")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Code int `json:"code"`
		Data struct {
			IsForbiddenRegion bool   `json:"IsForbiddenRegion"`
			CountryCode       string `json:"CountryCode"`
		} `json:"data"`
	}

	if err := json.Unmarshal(b, &res); err != nil {
		// Fallback for older block mechanism
		bodyStr := string(b)
		if strings.Contains(bodyStr, "Your region is not supported") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusUnexpected}
	}

	if res.Data.IsForbiddenRegion {
		return core.Result{Status: core.StatusNo}
	}

	if res.Data.CountryCode != "" {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.Data.CountryCode)}
	}

	return core.Result{Status: core.StatusUnexpected}
}
