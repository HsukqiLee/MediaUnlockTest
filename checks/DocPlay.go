package mediaunlocktest

import (
	"net/http"
)

func DocPlay(c http.Client) Result {
	resp, err := GET(c, "https://www.docplay.com/subscribe")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 303 {
		return Result{Status: StatusOK}
	}
	
	if resp.StatusCode == 307 {
		return Result{Status: StatusNo}
	}
	
	return Result{Status: StatusUnexpected}
}