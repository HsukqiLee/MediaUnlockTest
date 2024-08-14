package mediaunlocktest

import (
	"net/http"
)

func AMCPlus(c http.Client) Result {
	resp1, err := GET(c, "https://www.amcplus.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
    if resp1.StatusCode == 302 {
        resp2, err := GET(c, resp1.Header.Get("Location"))
	    if err != nil {
		    return Result{Status: StatusNetworkErr, Err: err}
	    }
	    defer resp2.Body.Close()
	    
	    if resp2.StatusCode == 301 && resp2.Header.Get("Location") == "https://www.amcplus.com/pages/geographic-restriction" {
	        return Result{Status: StatusNo}
	    }
	    return Result{Status: StatusUnexpected}
    }
    if resp1.StatusCode == 403 {
        return Result{Status: StatusBanned}
    }
    if resp1.StatusCode == 200 {
        return Result{Status: StatusOK}
    }
	return Result{Status: StatusUnexpected}
}