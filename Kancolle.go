package mediaunlocktest

import (
	"net/http"
)

func Kancolle(c http.Client) Result {
	resp, err := GET_Dalvik(c, "http://203.104.209.7/kcscontents/news/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
	
	if resp.StatusCode == 403 || resp.StatusCode == 302 {
		return Result{Status: StatusNo}
	}
	
	return Result{Status: StatusUnexpected}
}
