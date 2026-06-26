package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func WikipediaEditability(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://zh.wikipedia.org/w/index.php?title=Wikipedia%3A%E6%B2%99%E7%9B%92&action=edit")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString, "Banned") {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	if resp.StatusCode == 429 {
		return core.Result{Status: core.StatusBanned}
	}

	return core.Result{Status: core.StatusUnexpected}
}
