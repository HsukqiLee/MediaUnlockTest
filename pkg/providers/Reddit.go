package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func Reddit(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.reddit.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 429 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 200 || resp.StatusCode == 302 {
		return core.Result{Status: core.StatusOK}
	}

	if resp.StatusCode == 403 && strings.Contains(bodyString, "blocked") {
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}

