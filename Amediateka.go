package mediaunlocktest

import (
	"net/http"
)

func Amediateka(c http.Client) Result {
	resp, err := GET(c, "https://www.amediateka.ru/")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 301 && resp.Header.Get("Location") == "https://www.amediateka.ru/unavailable/index.html" {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}