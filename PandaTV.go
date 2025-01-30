package mediaunlocktest

import (
	"net/http"
)

func PandaTV(c http.Client) Result {
	resp, err := GET(c, "https://api.pandalive.co.kr/v1/live/play")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 400 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
