package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"errors"
	"strings"
)

func PrettyDerbyJP(c core.HttpClient) core.Result {
	resp, err := core.GETRaw(c, "https://api-umamusume.cygames.jp/", core.H{"User-Agent", core.UA_Dalvik})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "timeout") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
