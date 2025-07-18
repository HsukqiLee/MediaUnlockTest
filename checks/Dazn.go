package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func Dazn(c http.Client) Result {
	resp, err := PostJson(c, "https://startup.core.indazn.com/misl/v5/Startup",
		`{"LandingPageKey":"generic","Languages":"zh-CN,zh,en","Platform":"web","PlatformAttributes":{},"Manufacturer":"","PromoCode":"","Version":"2"}`,
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
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
		return Result{Status: StatusErr, Err: err}
	}
	if res.Region.IsAllowed {
		return Result{
			Status: StatusOK,
			Region: res.Region.GeolocatedCountry,
		}
	}
	return Result{
		Status: StatusNo,
		Info:   res.Region.DisallowedReason,
	}
}
