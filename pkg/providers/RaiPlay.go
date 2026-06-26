package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func RaiPlay(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://mediapolisvod.rai.it/relinker/relinkerServlet.htm?cont=VxXwi7UcqjApssSlashbjsAghviAeeqqEEqualeeqqEEqual&output=64")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString, "no_available") {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
