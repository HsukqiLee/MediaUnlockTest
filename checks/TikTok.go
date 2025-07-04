package mediaunlocktest

import (
	"io"
	"net/http"
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

func TikTok(c http.Client) Result {
	resp, err := GET(c, "https://www.tiktok.com/explore")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		return Result{Status: StatusFailed}
	}

	if strings.Contains(bodyString, "https://www.tiktok.com/hk/notfound") {
		return Result{Status: StatusNo, Region: "hk"}
	}

	if region := extractTikTokRegion(bodyString); region != "" {

		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}

	return Result{Status: StatusNo}
}
