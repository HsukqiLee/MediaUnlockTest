package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Spotify(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://spclient.wg.spotify.com/signup/public/v1/account",
		`birth_day=11&birth_month=11&birth_year=2000&collect_personal_info=undefined&creation_flow=&creation_point=https%3A%2F%2Fwww.spotify.com%2Fhk-en%2F&displayname=Gay%20Lord&gender=male&iagree=1&key=a1e486e2729f46d6bb368d6b2bcda326&platform=www&referrer=&send-email=0&thirdpartyemail=0&identifier_token=AgE6YTvEzkReHNfJpO114514`,
		core.H{"Accept-Language", "en"},
		core.H{"cache-control", "no-cache"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	var res struct {
		Status            int
		Country           string
		IsCountryLaunched bool `json:"is_country_launched"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Status == 320 {
		return core.Result{Status: core.StatusNo}
	}
	if res.Status == 311 && res.IsCountryLaunched {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(res.Country)}
	}
	return core.Result{Status: core.StatusNo}
}
