package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func Starz(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.starz.com/sapi/header/v1/starz/us/109448574b2147ccbc494b429ff5ef1b", core.H{"Referer", "https://www.starz.com/us/en/"})
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned, Err: err}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	authorization := string(b)
	resp2, err := core.GET(c, "https://auth.starz.com/api/v4/User/geolocation",
		core.H{"AuthTokenAuthorization", authorization},
		core.H{"BestAvailableToken", "true"},
		core.H{"Origin", "https://www.starz.com"},
		core.H{"Referer", "https://www.starz.com/"},
		core.H{"X-Client-Features", "DeviceCount"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		IsAllowedAccess  bool
		IsAllowedCountry bool
		IsKnownProxy     bool
		Country          string
	}
	if err := json.Unmarshal(b2, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.IsAllowedAccess && res.IsAllowedCountry && !res.IsKnownProxy {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
