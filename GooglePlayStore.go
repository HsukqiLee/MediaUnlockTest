package mediaunlocktest

import (
	//"fmt"
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

func GooglePlayStore(c http.Client) Result {
	resp, err := GET(c, "https://play.google.com/store/games")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode == 200 {
		region := extractGooglePlayStoreRegion(bodyString)
		if region != "" {
			if region == "CN" {
				return Result{Status: StatusNo, Region: "cn"}
			}
			return Result{Status: StatusOK, Region: strings.ToLower(region)}
		}
	}

	return Result{Status: StatusUnexpected}
}
