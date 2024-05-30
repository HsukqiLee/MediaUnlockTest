package mediaunlocktest

import (
	"net/http"
)

func ITVX(c http.Client) Result {
	resp, err := GET(c, "https://simulcast.itv.com/playlist/itvonline/ITV", H{"connection", "keep-alive"})
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 404  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}