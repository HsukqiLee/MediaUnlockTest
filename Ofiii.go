package mediaunlocktest

import (
	"net/http"
)

func Ofiii(c http.Client) Result {
	resp, err := GET(c, "https://ntdofifreepc.akamaized.net")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 451 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
