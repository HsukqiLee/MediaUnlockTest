package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	tls_client "github.com/bogdanfinn/tls-client"
	"regexp"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func extractBingRegion(responseBody string) string {
	re := regexp.MustCompile(`Region:"([^"]*)"`)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func Bing(c core.HttpClient) core.Result {
	jar := tls_client.NewCookieJar()
	c.SetCookieJar(jar)
	resp, err := core.GET(c, "https://www.bing.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.Header.Get("Location") == "https://www.bing.com/?brdr=1" {
		resp, err = core.GET(c, "https://www.bing.com/")
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp.Body.Close()
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if resp.StatusCode == 200 {
		region := extractBingRegion(bodyString)
		if region == "CN" {
			return core.Result{Status: core.StatusNo, Region: "cn"}
		}
		if region != "" {
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
		}
	}

	if strings.Contains(bodyString, "cn.bing.com") {
		return core.Result{Status: core.StatusNo, Region: "cn"}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
