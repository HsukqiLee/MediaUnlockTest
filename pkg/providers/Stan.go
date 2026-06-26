package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func Stan(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://api.stan.com.au/login/v1/sessions/web/account", `{}`)
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned, Err: err}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusUnexpected, Err: err}
	}

	if strings.Contains(bodyString, "Access Denied") {
		return core.Result{Status: core.StatusNo}
	}

	if strings.Contains(bodyString, "VPNDetected") {
		return core.Result{Status: core.StatusNo, Info: "VPN Detected"}
	}

	if resp.StatusCode == 400 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
