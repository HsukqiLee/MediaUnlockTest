package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

var deepseekRegionRegex = regexp.MustCompile(`<meta\s+name="region"\s+content="([^"]+)"`)

func DeepSeek(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://chat.deepseek.com/sign_in")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 403:
		return core.Result{Status: core.StatusNo}
	case 200:
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			matches := deepseekRegionRegex.FindStringSubmatch(string(b))
			if len(matches) > 1 {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
			}
		}
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
