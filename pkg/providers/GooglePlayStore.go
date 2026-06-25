package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func extractGooglePlayStoreRegion(responseBody string) string {
	re := regexp.MustCompile(`"zQmIje"\s*:\s*"([^"]+)"`)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func GooglePlayStore(c http.Client) core.Result {
	resp, err := core.GET(c, "https://play.google.com/store/games")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode == 200 {
		region := extractGooglePlayStoreRegion(bodyString)
		if region != "" {
			if region == "CN" {
				return core.Result{Status: core.StatusNo, Region: "cn"}
			}
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
		}
	}

	return core.Result{Status: core.StatusUnexpected}
}

