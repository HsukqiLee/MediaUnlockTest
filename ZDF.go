package mediaunlocktest

import (
	"net/http"
)

func ZDF(c http.Client) Result {
	resp, err :=  GET_Dalvik(c, "https://ssl.zdf.de/geo/de/geo.txt")
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