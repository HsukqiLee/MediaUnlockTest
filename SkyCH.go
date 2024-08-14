package mediaunlocktest

import (
    "net/http"
)

func Sky_CH(c http.Client) Result {
    resp, err := GET(c, "https://www.sky.ch/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	
	if resp.StatusCode == 302 && resp.Header.Get("Location") == "https://www.sky.ch/en/out-of-country" {
		return Result{Status: StatusNo}
	}
    
    if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
    
    return Result{Status: StatusUnexpected}
}
