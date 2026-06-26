package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func Dazn(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://startup.core.indazn.com/misl/v5/Startup",
		`{"LandingPageKey":"generic","Languages":"zh-CN,zh,en","Platform":"web","PlatformAttributes":{},"Manufacturer":"","PromoCode":"","Version":"2"}`,
	)
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
		Region struct {
			IsAllowed             bool   `json:"isAllowed"`
			DisallowedReason      string `json:"disallowedReason"`
			GeolocatedCountry     string `json:"geolocatedCountry"`
			GeolocatedCountryName string `json:"geolocatedCountryName"`
		} `json:"region"`
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Region.IsAllowed {
		return core.Result{
			Status: core.StatusOK,
			Region: res.Region.GeolocatedCountry,
		}
	}
	return core.Result{
		Status: core.StatusNo,
		Info:   res.Region.DisallowedReason,
	}
}
