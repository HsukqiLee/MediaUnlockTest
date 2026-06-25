package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func TNTSports(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.tntsports.co.uk/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 307 && resp.Header.Get("Location") == "https://www.tntsports.co.uk/geoblocking.shtml" {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}

		bodyString := string(body)

		re := regexp.MustCompile(`\\\"countryCode\\\":\\\"([A-Z]{2})\\\"`)
		matches2 := re.FindStringSubmatch(bodyString)
		if len(matches2) >= 2 {
			countryCode := matches2[1]
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(countryCode)}
		}
	}

	return core.Result{Status: core.StatusNo}
}

