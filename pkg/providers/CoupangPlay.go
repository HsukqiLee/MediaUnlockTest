package mediaunlocktest

import (
	"net/http"
)

func CoupangPlay(c http.Client) Result {
	resp, err := GET_Dalvik(c, "https://www.coupangplay.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 && resp.Header.Get("Location") == "https://www.coupangplay.com/not-available" {
		return Result{Status: StatusNo}
	}

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusBanned},
	}, Result{Status: StatusUnexpected})
}
