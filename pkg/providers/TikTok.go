package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

func extractTikTokRegion(body string) string {
	re := regexp.MustCompile(`"region":"(\w+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func TikTok(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.tiktok.com/explore")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString, "https://www.tiktok.com/hk/notfound") {
		return core.Result{Status: core.StatusNo, Region: "hk"}
	}

	if region := extractTikTokRegion(bodyString); region != "" {

		return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
	}

	return core.Result{Status: core.StatusNo}
}
