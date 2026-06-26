package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func VideoLand(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://api.videoland.com/subscribe/videoland-account/graphql", `{"operationName":"IsOnboardingGeoBlocked","variables":{},"query":"query IsOnboardingGeoBlocked {\n  isOnboardingGeoBlocked\n}\n"}`,
		core.H{"connection", "keep-alive"},
		core.H{"apollographql-client-name", "apollo_accounts_base"},
		core.H{"traceparent", "00-cab2dbd109bf1e003903ec43eb4c067d-623ef8e56174b85a-01"},
		core.H{"origin", "https://www.videoland.com"},
		core.H{"referer", "https://www.videoland.com/"},
		core.H{"accept", "application/json, text/plain, */*"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Data struct {
			Blocked bool `json:"isOnboardingGeoBlocked"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Data.Blocked {
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusOK}
}
