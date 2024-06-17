package mediaunlocktest

import (
	"net/http"
)

func MoviStarPlus(c http.Client) Result {
	resp, err := GET(c, "https://contratar.movistarplus.es/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}