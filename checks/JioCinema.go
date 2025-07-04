package mediaunlocktest

import (
	"net/http"
)

func JioCinema(c http.Client) Result {
	resp, err := GET(c, "https://content-jiovoot.voot.com/psapi/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 474:
		return Result{Status: StatusNo}
	case 200:
		return Result{Status: StatusOK}
	default:
		return Result{Status: StatusUnexpected}
	}
}
