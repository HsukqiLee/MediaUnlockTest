package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func Max(c core.HttpClient) core.Result {
	h1 := "beam/5.0.0 (desktop/desktop; Windows/10; afbb5daa-c327-461d-9460-d8e4b3ee4a1f/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)"
	h2 := "WEB:10:beam:5.2.1"
	h3 := "realm=bolt"

	resp1, err := core.GET(c, "https://default.any-any.prd.api.max.com/token?realm=bolt&deviceId=afbb5daa-c327-461d-9460-d8e4b3ee4a1f",
		core.H{"x-device-info", h1},
		core.H{"x-disco-client", h2},
		core.H{"x-disco-params", h3},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res1 struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}
	token := res1.Data.Attributes.Token

	resp2, err := core.PostJson(c, "https://default.any-any.prd.api.max.com/session-context/headwaiter/v1/bootstrap", `{}`,
		core.H{"Cookie", "st=" + token},
		core.H{"x-device-info", h1},
		core.H{"x-disco-client", h2},
		core.H{"x-disco-params", h3},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res2 struct {
		Routing struct {
			Domain     string `json:"domain"`
			Tenant     string `json:"tenant"`
			Env        string `json:"env"`
			HomeMarket string `json:"homeMarket"`
		} `json:"routing"`
	}
	if err := json.Unmarshal(body2, &res2); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	domain := res2.Routing.Domain
	tenant := res2.Routing.Tenant
	env := res2.Routing.Env
	homeMarket := res2.Routing.HomeMarket

	resp3, err := core.GET(c, "https://default."+tenant+"-"+homeMarket+"."+env+"."+domain+"/users/me",
		core.H{"Cookie", "st=" + token},
		core.H{"x-device-info", h1},
		core.H{"x-disco-client", h2},
		core.H{"x-disco-params", h3},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res3 struct {
		Data struct {
			Attributes struct {
				CurrentLocationTerritory string `json:"currentLocationTerritory"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body3, &res3); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}
	region := res3.Data.Attributes.CurrentLocationTerritory

	resp4, err := core.GET(c, "https://www.max.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	loc := resp4.Header.Get("Location")
	if (resp4.StatusCode >= 300 && resp4.StatusCode < 400 && loc != "" && (strings.Contains(loc, "hbomax.com") || strings.Contains(loc, "geo-availability"))) || region == "" {
		return core.Result{Status: core.StatusNo}
	}

	resp5, err := core.GET(c, "https://default.any-any.prd.api.max.com/any/playback/v1/playbackInfo",
		core.H{"Cookie", "st=" + token},
		core.H{"x-device-info", h1},
		core.H{"x-disco-client", h2},
		core.H{"x-disco-params", h3},
	)
	if err == nil {
		defer resp5.Body.Close()
		body5, _ := io.ReadAll(resp5.Body)
		if strings.Contains(string(body5), "VPN") {
			return core.Result{Status: core.StatusBanned}
		}
	}

	return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
}
