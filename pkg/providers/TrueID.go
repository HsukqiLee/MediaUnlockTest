package mediaunlocktest

import (
	"io"
	"net/http"
	"regexp"
	//"strings"
)

func extractTrueIDChannelID(body string) string {
	regex := regexp.MustCompile(`"channelId"\s*:\s*"([^"]+)`)
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractTrueIDAuthUser(body string) string {
	regex := regexp.MustCompile(`"buildId"\s*:\s*"([^"]+)`)
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractTrueIDBillboardType(body string) string {
	regex := regexp.MustCompile(`"billboardType"\s*:\s*"([^"]+)`)
	matches := regex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func TrueID(c http.Client) Result {
	resp1, err := GET(c, "https://tv.trueid.net/th-en/live/thairathtv-hd")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	channelId := extractTrueIDChannelID(string(body1))
	authUser := extractTrueIDAuthUser(string(body1))
	if len(authUser) < 11 {
		return Result{Status: StatusNo}
	}
	authKey := authUser[10:]

	req, err := http.NewRequest("GET", "https://tv.trueid.net/api/stream/checkedPlay?channelId="+channelId+"&lang=en&country=th", nil)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	req.Header.Set("user-agent", UA_Browser)
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.SetBasicAuth(authUser, authKey)
	resp2, err := cdo(c, req)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	switch extractTrueIDBillboardType(string(body2)) {
	case "GEO_BLOCK":
		return Result{Status: StatusNo}
	case "LOADING":
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
