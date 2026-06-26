package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"encoding/json"
	"io"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func MyTvSuper(c core.HttpClient) core.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "https://www.mytvsuper.com/api/auth/getSession/self/", nil)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	r.Header.Set("User-Agent", core.UA_Browser)
	r.Header.Set("Content-Type", "application/json")

	resp, err := core.Cdo(c, r)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res mytvsuperRes
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Region == 1 {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}

type mytvsuperRes struct {
	Region int
}
