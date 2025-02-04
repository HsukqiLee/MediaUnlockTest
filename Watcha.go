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
	resp1, err := GET(c, "https://watcha.com/api/aio_browses/tvod/all?size=3")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == 451 {
		return Result{Status: StatusNo}
	}

	resp2, err := GET(c, "https://watcha.com/browse/theater",
		H{"Connection", "keep-alive"},
		H{"Upgrade-Insecure-Requests", "1"},
		H{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 451 {
		return Result{Status: StatusNo}
	}
	if resp2.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	if resp2.StatusCode == 200 {
		return Result{Status: StatusOK, Region: "kr"}
	}
	if resp2.StatusCode == 302 {
		if resp2.Header.Get("Location") == "/ja-JP/browse/theater" {
			return Result{Status: StatusOK, Region: "jp"}
		}
		if resp2.Header.Get("Location") == "/ko-KR/browse/theater" {
			return Result{Status: StatusOK, Region: "kr"}
		}
	}
	return Result{Status: StatusUnexpected}
}
