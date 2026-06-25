package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func extractLineTVMainJsUrl(html string) string {
	re := regexp.MustCompile(`href="(https://web-static\.linetv\.tw/release-fargate/public/dist/main-[a-z0-9]{8}-prod\.js)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractLineTVAppId(js string) string {
	re := regexp.MustCompile(`appId:"([^"]+)"`)
	matches := re.FindStringSubmatch(js)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func LineTV(c http.Client) core.Result {
	resp1, err := core.GET(c, "https://www.linetv.tw/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString1 := string(body1)
	mainJsUrl := extractLineTVMainJsUrl(bodyString1)

	if mainJsUrl == "" {
	    return core.Result{Status: core.StatusFailed}
	}
	resp2, err := core.GET(c, mainJsUrl)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString2 := string(body2)
	appId := extractLineTVAppId(bodyString2)
	if appId == "" {
	    return core.Result{Status: core.StatusFailed}
	}
	
	resp3, err := core.GET(c, "https://www.linetv.tw/api/part/11829/eps/1/part?appId=" + appId + "&productType=FAST&version=10.38.0")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()
	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		CountryCode int `json:"countryCode"`
	}
	if err := json.Unmarshal(body3, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.CountryCode == 228 {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}



