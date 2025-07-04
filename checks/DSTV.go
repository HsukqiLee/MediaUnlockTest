package mediaunlocktest

import (
	"net/http"
)

func DSTV(c http.Client) Result {
	resp, err := GET(c, "https://now.dstv.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 451:
		return Result{Status: StatusNo}
	case 200:
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
