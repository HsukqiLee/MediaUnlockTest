package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func AnimeFesta(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api-animefesta.iowl.jp/v1/titles/1560",
		core.H{"accept", "application/json"},
		core.H{"accept-language", "en-US,en;q=0.9"},
		core.H{"anime-user-tracking-id", "yEZr4P_U7JEdBucZOkv1Y"},
		core.H{"authorization", ""},
		core.H{"origin", "https://animefesta.iowl.jp"},
		core.H{"referer", "https://animefesta.iowl.jp/"},
		core.H{"sec-gpc", "1"},
		core.H{"x-requested-with", "XMLHttpRequest"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

