package mediaunlocktest

import (
	"net/http"
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
	return Result{Status: StatusNetworkErr, Err: err}
}
bodyString := string(bodyBytes)
region := extractWatchaRegion(bodyString)
if region != "" {
	return Result{Status: StatusOK, Region: strings.ToLower(region)}
}
*/

func Watcha(c http.Client) Result {
	resp, err := GET(c, "https://watcha.com/browse/theater",
		H{"Connection", "keep-alive"},
		H{"Upgrade-Insecure-Requests", "1"},
		H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 451 {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK, Region: "kr"}
	}
	if resp.StatusCode == 302 {
		if resp.Header.Get("Location") == "https://watcha.com/jp-ja" {
			return Result{Status: StatusOK, Region: "jp"}
		}
		if resp.Header.Get("Location") == "/ko-KR/browse/theater" {
			return Result{Status: StatusOK, Region: "kr"}
		}
	}
	return Result{Status: StatusUnexpected}
}
