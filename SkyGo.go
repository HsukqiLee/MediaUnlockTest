package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func SkyGo(c http.Client) Result {
	resp, err := GET(c, "https://skyid.sky.com/authorise/skygo?response_type=token&client_id=sky&appearance=compact&redirect_uri=skygo://auth")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if resp.StatusCode == 302 {
		return Result{Status: StatusOK}
	}
	
	if resp.StatusCode == 403 || strings.Contains(bodyString, "Access Denied") {
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
