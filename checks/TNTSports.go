package mediaunlocktest

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

func TNTSports(c http.Client) Result {
	resp, err := GET(c, "https://www.tntsports.co.uk/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 307 && resp.Header.Get("Location") == "https://www.tntsports.co.uk/geoblocking.shtml" {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		bodyString := string(body)

		re := regexp.MustCompile(`\\\"countryCode\\\":\\\"([A-Z]{2})\\\"`)
		matches2 := re.FindStringSubmatch(bodyString)
		if len(matches2) >= 2 {
			countryCode := matches2[1]
			return Result{Status: StatusOK, Region: strings.ToLower(countryCode)}
		}
	}

	return Result{Status: StatusUnexpected}
}
