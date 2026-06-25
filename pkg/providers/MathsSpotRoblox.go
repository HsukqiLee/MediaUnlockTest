package providers

import (
	core "MediaUnlockTest/pkg/core"

	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func extractMathsSpotRobloxAPIPath(body string) string {
	re := regexp.MustCompile(`fetch\("(/[^"]+)\/reportEvent"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractMathsSpotRobloxRegion(body string) string {
	re := regexp.MustCompile(`"countryCode"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractMathsSpotRobloxNggFeVersion(body string) string {
	re := regexp.MustCompile(`"NEXT_PUBLIC_FE_VERSION"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func MathsSpotRoblox(c http.Client) core.Result {
	resp1, err := core.GET(c, "https://mathsspot.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	bodyString1 := string(body1)
	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString1, "FailureServiceNotInRegion") {
		return core.Result{Status: core.StatusNo}
	}

	apiPath := extractMathsSpotRobloxAPIPath(bodyString1)
	region := extractMathsSpotRobloxRegion(bodyString1)
	nggFeVersion := extractMathsSpotRobloxNggFeVersion(bodyString1)
	if nggFeVersion == "berlin-v1.34.800_redisexp-arm.1" {
		nggFeVersion = "berlin-v1.34.810_redisexp-arm.1"
	}

	if apiPath == "" || region == "" || nggFeVersion == "" {
		return core.Result{Status: core.StatusFailed}
	}

	fakeUAId := core.GenRandomStr(21)
	fakeSessId := core.GenRandomStr(21)
	fakeFesessId := core.GenRandomStr(21)
	fakeVisitId := core.GenRandomStr(21)

	resp2, err := core.GET(c, "https://mathsspot.com"+apiPath+"/startSession?appId=5349&uaId=ua-"+fakeUAId+"&uaSessionId=uasess-"+fakeSessId+"&feSessionId=fesess-"+fakeFesessId+"&visitId=visitid-"+fakeVisitId+"&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%%3A%%2F%%2Fmathsspot.com%%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA",
		core.H{"referer", "https://mathsspot.com/"},
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
		nggFeVersion = string(body2)
		resp3, err := core.GET(c, "https://mathsspot.com"+apiPath+"/startSession?appId=5349&uaId=ua-"+fakeUAId+"&uaSessionId=uasess-"+fakeSessId+"&feSessionId=fesess-"+fakeFesessId+"&visitId=visitid-"+fakeVisitId+"&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%%3A%%2F%%2Fmathsspot.com%%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA",
			core.H{"referer", "https://mathsspot.com/"},
			core.H{"x-ngg-skip-evar-check", "true"},
			core.H{"x-ngg-fe-version", nggFeVersion},
		)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()

		body3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		if err := json.Unmarshal(body3, &res); err != nil {
			return core.Result{Status: core.StatusFailed, Err: err}
		}
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
