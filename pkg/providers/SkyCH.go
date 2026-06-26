package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
)

func Sky_CH(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://gateway.prd.sky.ch/user/customer/create")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusErr, Err: err}
		}
		if string(bodyBytes) == `{"message": "", "code": "GEO_BLOCKED"}` {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 405 {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
