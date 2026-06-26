package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func MXPlayer(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.mxplayer.in/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}

		bodyString := string(body)
		if strings.Contains(bodyString, "We are currently not available in your region") {
			return core.Result{Status: core.StatusNo}
		}

		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
