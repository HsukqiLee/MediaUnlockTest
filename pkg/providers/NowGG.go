package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

func NowGG(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://now.gg/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}
	bodyString := string(body)

	if strings.Contains(bodyString, "FailureServiceNotInRegion") {
		return core.Result{Status: core.StatusNo}
	}

	re := regexp.MustCompile(`"countryCode"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(bodyString)
	if len(matches) > 1 {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
	}

	return core.Result{Status: core.StatusFailed}
}
