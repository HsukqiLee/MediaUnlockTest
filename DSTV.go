package mediaunlocktest

import (
	"net/http"
)

func DSTV(c http.Client) Result {
	resp, err := GET(c, "https://authentication.dstv.com/favicon.ico")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 404  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}