package providers

import (
	"MediaUnlockTest/pkg/core"
	"regexp"
	"strings"
)

func Zee5(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.zee5.com/global")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	cookies := resp.Header.Values("Set-Cookie")
	re := regexp.MustCompile(`country=([A-Z]{2})`)

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}

	for _, cookie := range cookies {
		matches := re.FindStringSubmatch(cookie)
		if len(matches) > 1 {
			region := matches[1]
			if region != "" {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
			}
		}
	}

	return core.Result{Status: core.StatusUnexpected}
}
