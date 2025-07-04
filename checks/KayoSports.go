package mediaunlocktest

import (
	"net/http"
)

func KayoSports(c http.Client) Result {
	resp, err := GET(c, "https://kayosports.com.au",
	    H{"Accept", "*/*"},
	    H{"Accept-Language", "en-US,en;q=0.9"},
	    H{"Origin", "https://kayosports.com.au"},
	    H{"Referer", "https://kayosports.com.au/"},
	)
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