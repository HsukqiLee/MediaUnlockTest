package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func SlingTV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://p-geo.movetv.com/geo")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		IpRestricted bool   `json:"ip_restricted"`
		Country      string `json:"country"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if res.IpRestricted {
		return core.Result{Status: core.StatusNo, Region: strings.ToLower(res.Country)}
	}
	return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.Country)}
}
