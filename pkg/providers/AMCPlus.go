package providers

import (
	"MediaUnlockTest/pkg/core"
	"regexp"

	http "github.com/bogdanfinn/fhttp"
)

func extractAMCPlusRegion(url string) string {
	re := regexp.MustCompile(`https://www\.amcplus\.com/countries/(\w{2})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func AMCPlus(c core.HttpClient) core.Result {
	resp1, err := core.GET(c, "https://www.amcplus.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	if resp1.StatusCode == 302 {
		resp2, err := core.GET(c, resp1.Header.Get("Location"))
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()

		if resp2.StatusCode == 301 {
			if resp2.Header.Get("Location") == "https://www.amcplus.com/pages/geographic-restriction" {
				return core.Result{Status: core.StatusNo}
			}
			resp3, err := core.GET(c, resp2.Header.Get("Location"))
			if err != nil {
				return core.Result{Status: core.StatusNetworkErr, Err: err}
			}
			defer resp3.Body.Close()
			if resp3.StatusCode == 200 {
				region := extractAMCPlusRegion(resp1.Header.Get("Location"))
				if region != "" {
					return core.Result{Status: core.StatusOK, Region: region}
				}
			}
		}
		return core.Result{Status: core.StatusUnexpected}
	}

	return core.ResultFromMapping(resp1.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK, Region: "us"},
		http.StatusForbidden: {Status: core.StatusBanned},
	}, core.Result{Status: core.StatusUnexpected})
}
