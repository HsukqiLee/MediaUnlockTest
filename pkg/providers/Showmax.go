package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"regexp"
	"strings"
)

func extractShowmaxRegion(body string) string {
	re := regexp.MustCompile(`hterr=([A-Z]{2})`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Showmax(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.showmax.com/",
		core.H{"Connection", "keep-alive"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return core.Result{Status: core.StatusNo}
	}

	region := extractShowmaxRegion(cookie)
	if region != "" {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
	}
	return core.Result{Status: core.StatusUnexpected}
}

