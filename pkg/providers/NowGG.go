package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func extractNowGGAPIPath(body string) string {
	re := regexp.MustCompile(`fetch\("(/[^"]+)/oapi/1/play/v1/push"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractNowGGRegion(body string) string {
	re := regexp.MustCompile(`"countryCode"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractNowGGNggFeVersion(body string) string {
	re := regexp.MustCompile(`"NEXT_PUBLIC_FE_VERSION"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func NowGG(c http.Client) core.Result {
	resp1, err := core.GET(c, "https://now.gg/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}
	bodyString1 := string(body1)

	if strings.Contains(bodyString1, "FailureServiceNotInRegion") {
		return core.Result{Status: core.StatusNo}
	}

	apiPath := extractNowGGAPIPath(bodyString1)
	region := extractNowGGRegion(bodyString1)
	nggFeVersion := extractNowGGNggFeVersion(bodyString1)

	if apiPath == "" || region == "" || nggFeVersion == "" {
		return core.Result{Status: core.StatusFailed}
	}

	fakeUAId := core.GenRandomStr(21)
	fakeSessId := core.GenRandomStr(21)
	fakeFesessId := core.GenRandomStr(21)
	fakeVisitId := core.GenRandomStr(21)

	baseUrl := "https://now.gg" + apiPath + "/startSession?appId=5349&uaId=ua-" + fakeUAId + "&uaSessionId=uasess-" + fakeSessId + "&feSessionId=fesess-" + fakeFesessId + "&visitId=visitid-" + fakeVisitId + "&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%3A%2F%2Fnow.gg%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA"

	resp2, err := core.GET(c, baseUrl,
		core.H{"referer", "https://now.gg/"},
		core.H{"x-ngg-skip-evar-check", "true"},
		core.H{"x-ngg-fe-version", nggFeVersion},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body2, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	switch res.Status {
	case "FailureServiceNotInRegion":
		return core.Result{Status: core.StatusNo}
	case "FailureProxyUserLimitExceeded":
		return core.Result{Status: core.StatusNo, Info: "Proxy Detected"}
	case "Success", "FailureUnauthorized":
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
	}

	return core.Result{Status: core.StatusFailed}
}
