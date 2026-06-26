package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func Abema(c core.HttpClient) core.Result {
	resp, err := core.GET_Dalvik(c, "https://api.abema.io/v1/ip/check?device=android")
	if err != nil {
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
		IsoCountryCode string
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.IsoCountryCode == "JP" {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusRestricted, Info: "Oversea Only"}
}
