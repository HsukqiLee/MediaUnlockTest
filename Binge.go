package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func Binge(c http.Client) Result {
	resp, err := GET(c, "https://auth.streamotion.com.au")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if strings.Contains(bodyString, "Access Denied") {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}