package mediaunlocktest

import (
	"net/http"
)

func  AnimeFesta(c http.Client) Result {
	resp, err := GET(c, "https://api-animefesta.iowl.jp/v1/titles/1305")
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