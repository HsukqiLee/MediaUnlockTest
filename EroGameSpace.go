package mediaunlocktest

import (
    "io"
	"net/http"
	"strings"
)

func EroGameSpace(c http.Client) Result {
	resp, err := GET(c, "https://erogamescape.org/~ap2/ero/toukei_kaiseki/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 200 {
	    if strings.Contains(bodyString, "18歳") {
		    return Result{Status: StatusOK}
	    }
	    return Result{Status: StatusNo}
	}
	
	return Result{Status: StatusUnexpected}
}