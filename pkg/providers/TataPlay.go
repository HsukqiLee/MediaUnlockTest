package providers

import (
	"MediaUnlockTest/pkg/core"
)

func TataPlay(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://watch.tataplay.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return core.Result{Status: core.StatusOK}
	case 403:
		return core.Result{Status: core.StatusNo}
	default:
		return core.Result{Status: core.StatusUnexpected}
	}
}
