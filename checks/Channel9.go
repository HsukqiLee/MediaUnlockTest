package mediaunlocktest

import (
	"net/http"
)

func Channel9(c http.Client) Result {
	resp, err := GET(c, "https://login.nine.com.au")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 302  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}