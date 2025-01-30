package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func HBOMax(c http.Client) Result {
	resp1, err := GET(c, "https://default.any-any.prd.api.max.com/token?realm=bolt&deviceId=afbb5daa-c327-461d-9460-d8e4b3ee4a1f",
		H{"x-device-info", "beam/5.0.0 (desktop/desktop; Windows/10; afbb5daa-c327-461d-9460-d8e4b3ee4a1f/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)"},
		H{"x-disco-client", "WEB:10:beam:5.2.1"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res1 struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	token := res1.Data.Attributes.Token

	resp2, err := PostJson(c, "https://default.any-any.prd.api.max.com/session-context/headwaiter/v1/bootstrap", ``,
		H{"Cookie", "st=" + token},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
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
		return Result{Status: StatusFailed, Err: err}
	}

	domain := res2.Routing.Domain
	tenant := res2.Routing.Tenant
	env := res2.Routing.Env
	homeMarket := res2.Routing.HomeMarket

	resp3, err := GET(c, "https://default."+tenant+"-"+homeMarket+"."+env+"."+domain+"/users/me",
		H{"Cookie", "st=" + token},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res3 struct {
		Data struct {
			Attributes struct {
				CurrentLocationTerritory string `json:"currentLocationTerritory"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body3, &res3); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	region := res3.Data.Attributes.CurrentLocationTerritory

	resp4, err := GET(c, "https://www.max.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp4.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	availableRegion := strings.ToUpper(strings.Join(strings.Fields(strings.Join(strings.Split(string(body4), "\"url\":\"/")[1:], " ")), " "))

	if strings.Contains(availableRegion, region) && region != "" {
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}
	return Result{Status: StatusNo}
}
