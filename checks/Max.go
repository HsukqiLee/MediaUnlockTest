package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Max(c http.Client) Result {
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
		resp5, err := PostForm(c, "https://default.any-any.prd.api.max.com/any/playback/v1/playbackInfo", `st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi0wOWQxOTg4Yy1mZmUzLTQxMDEtOWI5My0yNDU1ZTkyNGQ1YjYiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6YjYzOTgxZWQtNzA2MC00ZGYwLThkZGItZjA2YjFkNWRjZWVkIiwiaWF0IjoxNzQzODQwMzgwLCJleHAiOjIwNTkyMDAzODAsInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fYW1lciIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6IjQwYTgzZjNlLTY4OTktNDE3Mi1hMWY2LWJjZDVjN2ZkNjA4NSIsInZlcnNpb24iOiJ2MyIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiNWY3YzViZjQtYjc4Ny00NDRjLWJhYTYtMzU5MzgwYWFiM2RmIn0.f5HTgIV2v0nQQDp5LQG0xqLrxyACdvnMDiWO_viX_CUGqtc5ncSjp_LgM30QFkkMnINFhzKEGRpsZvb-o3Pj_Z39uRBr5LCeiCPR7ssV-_SXyRFVRRDEB2lpxyz7jmdD1SxvA06HnEwTbZQzlbZ7g9GXq02yNdEfHlqYEh_4WF88UbXfeieYTd4TH7kwN1RE50NfQUS6f0WmzpAbpiULyd87mpTeynchFNMMz-YHVzZ_-nDW6geihXc3tS0FKVSR8fdOSPQFzEYOLCfhInufiPahiXI-OKF89aShAqM-y4Hx_eukGnsq3mO5wa3unnqVr9Kzc61BIhHh1Hs2bqYiYg`)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		body5, err := io.ReadAll(resp5.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(body5), "VPN") {
			return Result{Status: StatusNo, Info: "VPN Detected"}
		}
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}
	return Result{Status: StatusNo}
}
