package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"regexp"
)

func extractSkyShowTimeRegion(url string) string {
	re := regexp.MustCompile(`https://www.skyshowtime.com/([a-z]{2})\?`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func SkyShowTime(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.skyshowtime.com/",
		core.H{"Cookie", "sat_track=true; AMCVS_99B971AC61C1E36F0A495FC6@AdobeOrg=1; AMCV_99B971AC61C1E36F0A495FC6@AdobeOrg=179643557|MCIDTS|19874|MCMID|36802229575946481753961418923958457479|MCOPTOUT-1717079521s|NONE|vVersion|5.5.0"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 307 {
		if resp.Header.Get("Location") == "https://www.skyshowtime.com/where-can-i-stream" {
			return core.Result{Status: core.StatusNo}
		}
		region := extractSkyShowTimeRegion(resp.Header.Get("Location"))
		if region != "" {
			return core.Result{Status: core.StatusOK, Region: region}
		}
		return core.Result{Status: core.StatusFailed}
	}

	return core.Result{Status: core.StatusUnexpected}
}

