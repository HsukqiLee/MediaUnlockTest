package mediaunlocktest

import (
    "net/http"
)

func ThreeNow(c http.Client) Result {
	resp, err := GET(c, "https://bravo-livestream.fullscreen.nz/index.m3u8")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}