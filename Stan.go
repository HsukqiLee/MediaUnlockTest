package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func Stan(c http.Client) Result {
	resp, err := PostJson(c, "https://api.stan.com.au/login/v1/sessions/web/account", `{}`)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if strings.Contains(bodyString, "Access Denied") {
		return Result{Status: StatusNo, Info: "Unavailable"}
	}
	
	if strings.Contains(bodyString, "VPNDetected") {
		return Result{Status: StatusNo, Info: "VPN Detected"}
	}
	
	if resp.StatusCode == 400 {
		return Result{Status: StatusOK}
	}
    
	return Result{Status: StatusUnexpected}
}