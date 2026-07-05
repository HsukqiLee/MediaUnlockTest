package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func DMM(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://api.tv.dmm.com/graphql", `{"query":"query FetchClient { client { isForeignAccess } }"}`)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, `"isForeignAccess":true`) {
		return core.Result{Status: core.StatusNo}
	}
	if strings.Contains(s, `"isForeignAccess":false`) {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
