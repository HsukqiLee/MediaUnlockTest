package mediaunlocktest

import (
	"net/http"
	"io"
	"encoding/json"
)

func VideoLand(c http.Client) Result {
	resp, err := PostJson(c, "https://api.videoland.com/subscribe/videoland-account/graphql", `{"operationName":"IsOnboardingGeoBlocked","variables":{},"query":"query IsOnboardingGeoBlocked {\n  isOnboardingGeoBlocked\n}\n"}`,
	    H{"connection", "keep-alive"},
	    H{"apollographql-client-name", "apollo_accounts_base"},
	    H{"traceparent", "00-cab2dbd109bf1e003903ec43eb4c067d-623ef8e56174b85a-01"},
	    H{"origin", "https://www.videoland.com"},
	    H{"referer", "https://www.videoland.com/"},
	    H{"accept", "application/json, text/plain, */*"},
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
		Data struct {
		    Blocked bool `json:"isOnboardingGeoBlocked"`
        } `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}
	if res.Data.Blocked {
		return Result{Status: StatusNo}
	}
    
	return Result{Status: StatusOK}
}