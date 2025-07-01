package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Copilot(c http.Client) Result {
	resp, err := GET(c, "https://copilot.microsoft.com/c/api/user?api-version=2")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		var res struct {
			RegionCode string `json:"regionCode"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return Result{Status: StatusErr, Err: err}
		}
		if res.RegionCode != "" {
			return Result{Status: StatusOK, Region: strings.ToLower(res.RegionCode)}
		}
	}
	return Result{Status: StatusUnexpected}
}
