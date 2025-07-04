package mediaunlocktest

import (
	"net/http"
	"regexp"
	"strings"
)

func Zee5(c http.Client) Result {
	resp, err := GET(c, "https://www.zee5.com/global")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	cookies := resp.Header.Values("Set-Cookie")
	re := regexp.MustCompile(`country=([A-Z]{2})`)

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	for _, cookie := range cookies {
		matches := re.FindStringSubmatch(cookie)
		if len(matches) > 1 {
			region := matches[1]
			if region != "" {
				return Result{Status: StatusOK, Region: strings.ToLower(region)}
			}
		}
	}

	return Result{Status: StatusUnexpected}
}
