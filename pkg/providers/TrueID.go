package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/base64"
	"io"
	"regexp"
	// "strings"
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

func TrueID(c core.HttpClient) core.Result {
	resp1, err := core.GET(c, "https://tv.trueid.net/th-en/live/thairathtv-hd")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	channelId := extractTrueIDChannelID(string(body1))
	authUser := extractTrueIDAuthUser(string(body1))
	if len(authUser) < 11 {
		return core.Result{Status: core.StatusNo}
	}
	authKey := authUser[10:]

	auth := base64.StdEncoding.EncodeToString([]byte(authUser + ":" + authKey))
	resp2, err := core.GET(c, "https://tv.trueid.net/api/stream/checkedPlay?channelId="+channelId+"&lang=en&country=th",
		core.H{"Authorization", "Basic " + auth},
	)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	switch extractTrueIDBillboardType(string(body2)) {
	case "GEO_BLOCK":
		return core.Result{Status: core.StatusNo}
	case "LOADING":
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
