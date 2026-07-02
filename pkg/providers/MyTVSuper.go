package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func MyTvSuper(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.mytvsuper.com/api/auth/getSession/self/", core.H{"Content-Type", "application/json"})
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
