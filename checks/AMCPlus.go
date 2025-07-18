package mediaunlocktest

import (
	"net/http"
	"regexp"
)

func extractAMCPlusRegion(url string) string {
	re := regexp.MustCompile(`https://www\.amcplus\.com/countries/(\w{2})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func AMCPlus(c http.Client) Result {
	resp1, err := GET(c, "https://www.amcplus.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	if resp1.StatusCode == 302 {
		resp2, err := GET(c, resp1.Header.Get("Location"))
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()

		if resp2.StatusCode == 301 {
			if resp2.Header.Get("Location") == "https://www.amcplus.com/pages/geographic-restriction" {
				return Result{Status: StatusNo}
			}
			resp3, err := GET(c, resp2.Header.Get("Location"))
			if err != nil {
				return Result{Status: StatusNetworkErr, Err: err}
			}
			defer resp3.Body.Close()
			if resp3.StatusCode == 200 {
				region := extractAMCPlusRegion(resp1.Header.Get("Location"))
				if region != "" {
					return Result{Status: StatusOK, Region: region}
				}
			}
		}
		return Result{Status: StatusUnexpected}
	}

	return ResultFromMapping(resp1.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK, Region: "us"},
		http.StatusForbidden: {Status: StatusBanned},
	}, Result{Status: StatusUnexpected})
}
