package providers

import (
	"MediaUnlockTest/pkg/core"
	//"regexp"
	//"strings"
)

/*
 * 如果以下方法失效，则使用这个方法
func extractWatchaRegion(responseBody string) string {
	re := regexp.MustCompile(`"regionCode"\s*:\s*"([^"]+)"`)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

bodyBytes, err := io.ReadAll(resp.Body)
if err != nil {
	return core.Result{Status: core.StatusNetworkErr, Err: err}
}
bodyString := string(bodyBytes)
region := extractWatchaRegion(bodyString)
if region != "" {
	return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
}
*/

func Watcha(c core.HttpClient) core.Result {
	resp1, err := core.GET(c, "https://watcha.com/api/aio_browses/tvod/all?size=3")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == 451 {
		return core.Result{Status: core.StatusNo}
	}

	resp2, err := core.GET(c, "https://watcha.com/browse/theater",
		core.H{"Connection", "keep-alive"},
		core.H{"Upgrade-Insecure-Requests", "1"},
		core.H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 451 {
		return core.Result{Status: core.StatusNo}
	}
	if resp2.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}
	if resp2.StatusCode == 200 {
		return core.Result{Status: core.StatusOK, Region: "kr"}
	}
	if resp2.StatusCode == 302 {
		if resp2.Header.Get("Location") == "/ja-JP/browse/theater" {
			return core.Result{Status: core.StatusOK, Region: "jp"}
		}
		if resp2.Header.Get("Location") == "/ko-KR/browse/theater" {
			return core.Result{Status: core.StatusOK, Region: "kr"}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}
