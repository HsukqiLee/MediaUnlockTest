package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"net/http"
	"time"
)

func KonosubaFD(c http.Client) core.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.konosubafd.jp/api/masterlist", nil)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	req.Header.Set("User-Agent", "pj0007/212 CFNetwork/1240.0.4 Darwin/20.6.0")

	resp, err := core.Cdo(c, req)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

