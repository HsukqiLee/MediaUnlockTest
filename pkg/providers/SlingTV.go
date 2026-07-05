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

	region := strings.ToLower(core.ThreeToTwoCode(res.Country))
	if res.IpRestricted {
		return core.Result{Status: core.StatusNo, Region: region}
	}
	return core.Result{Status: core.StatusOK, Region: region}
}
