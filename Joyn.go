package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func Joyn(c http.Client) Result {
	url := "https://auth.joyn.de/auth/anonymous"
	resp, err := PostJson(c, url,
		`{"client_id":"b74b9f27-a994-4c45-b7eb-5b81b1c856e7","client_name":"web","anon_device_id":"b74b9f27-a994-4c45-b7eb-5b81b1c856e7"}`,
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		// log.Println(err)
		return Result{Status: StatusFailed, Err: err}
	}
	url2 := "https://api.joyn.de/content/entitlement-token"
	resp2, err := PostJson(c, url2, `{"content_id":"daserste-de-hd","content_type":"LIVE"}`,
	    H{"authorization", "Bearer "+res.AccessToken},
	    H{"x-api-key", "36lp1t4wto5uu2i2nk57ywy9on1ns5yg"},
	)
	if err != nil {
		// log.Println(err)
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		// log.Println(err)
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res2a []struct {
        Code string `json:"code"`
    }
    var res2b struct {
        Token string `json:"entitlement_token"`
    }
    
	if err := json.Unmarshal(b2, &res2a); err != nil {
	    if err := json.Unmarshal(b2, &res2b); err != nil {
	       return Result{Status: StatusFailed, Err: err}
	    }
	    if res2b.Token != "" {
	        return Result{Status: StatusOK}
    	}
	}
	if res2a[0].Code == "ENT_AssetNotAvailableInCountry" {
	    return Result{Status: StatusNo}
	}
	return Result{Status: StatusUnexpected}
}
