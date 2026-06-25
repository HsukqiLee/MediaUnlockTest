package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Sooka(c http.Client) core.Result {
	resp, err := core.GET(c, "https://sooka.my/")
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

