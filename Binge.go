package mediaunlocktest

import (
	"net/http"
)

func Binge(c http.Client) Result {
	resp, err := GET(c, "https://auth.streamotion.com.au")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 302 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}