package mediaunlocktest

import (
	"net/http"
)

func Crackle(c http.Client) Result {
	resp, err := GET(c, "https://prod-api.crackle.com/appconfig",
        H{"Origin", "https://www.crackle.com"},
        H{"Referer", "https://www.crackle.com/"},
        H{"x-crackle-apiversion", "v2.0.0"},
        H{"x-crackle-brand", "crackle"},
        H{"x-crackle-platform", "5FE67CCA-069A-42C6-A20F-4B47A8054D46"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.Header.Get("x-crackle-region") == "US" {
	    return Result{Status: StatusOK}
	}
	
	if resp.Header.Get("x-crackle-region") != "" {
	    return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
