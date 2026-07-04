package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"errors"
	"strings"
)

func PrettyDerbyJP(c core.HttpClient) core.Result {
	resp, err := core.GETRaw(c, "https://api.games.umamusume.jp/", core.H{"User-Agent", core.UA_Dalvik})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		403: {Status: core.StatusNo},
		404: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
