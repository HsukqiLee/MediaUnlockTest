package mediaunlocktest

import (
	"io"
	"net/http"
)

func Sky_CH(c http.Client) Result {
	resp, err := GET(c, "https://gateway.prd.sky.ch/user/customer/create")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{Status: StatusErr, Err: err}
		}
		if string(bodyBytes) == `{"message": "", "code": "GEO_BLOCKED"}` {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusBanned}
	}

	if resp.StatusCode == 405 {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
