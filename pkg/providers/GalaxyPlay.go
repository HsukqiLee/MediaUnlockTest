package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func GalaxyPlay(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.glxplay.io/account/device/new")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return core.Result{Status: core.StatusOK}
	}

	if resp.StatusCode == 400 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		bodyStr := string(body)
		if len(bodyStr) > 0 && bodyStr[0] == '<' {
			return core.Result{Status: core.StatusNo}
		}
		
		if strings.Contains(bodyStr, `"errorCode": 495`) || strings.Contains(bodyStr, "not available in your region") {
			return core.Result{Status: core.StatusNo}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}
