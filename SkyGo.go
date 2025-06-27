package mediaunlocktest

import (
	"net/http"
)

func SkyGo(c http.Client) Result {
	resp, err := GET(c, "https://skyid.sky.com/authorise/skygo?response_type=token&client_id=sky&appearance=compact&redirect_uri=skygo://auth")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		return Result{Status: StatusOK}
	}

	if resp.StatusCode == 403 || resp.StatusCode == 200 {
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}

func SkyGo_NZ(c http.Client) Result {
	resp, err := GET(c, "https://linear-s.stream.skyone.co.nz/sky-sport-1.mpd")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return Result{Status: StatusOK}
	case 403:
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
