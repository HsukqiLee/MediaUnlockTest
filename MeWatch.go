package mediaunlocktest

import (
	"net/http"
)

func MeWatch(c http.Client) Result {
	resp, err := GET(c, "https://cdn.mewatch.sg/api/items/97098/videos?delivery=stream%2Cprogressive&ff=idp%2Cldp%2Crpt%2Ccd&lang=en&resolution=External&segments=all")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}