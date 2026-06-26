package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func ClipTV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://cliptv.vn/truyen-hinh")
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
		if strings.Contains(bodyString, "Sorry, this video is not available in your country.") {
			return core.Result{Status: core.StatusNo}
		}

		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
