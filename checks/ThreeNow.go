package mediaunlocktest

import (
	"net/http"
)

func ThreeNow(c http.Client) Result {
	resp, err := GET(c, "https://bravo-livestream.fullscreen.nz/index.m3u8")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
