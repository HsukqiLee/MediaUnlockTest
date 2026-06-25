package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func SkyGo(c http.Client) core.Result {
	resp, err := core.GET(c, "https://skyid.sky.com/authorise/skygo?response_type=token&client_id=sky&appearance=compact&redirect_uri=skygo://auth")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		return core.Result{Status: core.StatusOK}
	}

	if resp.StatusCode == 403 || resp.StatusCode == 200 {
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}

func SkyGo_NZ(c http.Client) core.Result {
	resp, err := core.GET(c, "https://linear-s.stream.skyone.co.nz/sky-sport-1.mpd")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return core.Result{Status: core.StatusOK}
	case 403:
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}

