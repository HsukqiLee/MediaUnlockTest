package mediaunlocktest

import (
	"net/http"
)

func CanalPlus(c http.Client) Result {
	resp, err := GET(c, "https://boutique-tunnel.canalplus.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 302 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}
	
	if resp.StatusCode == 403  {
		return Result{Status: StatusBanned}
	}

	return Result{Status: StatusUnexpected}
}