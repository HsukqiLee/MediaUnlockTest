package mediaunlocktest

import (
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

func MathsSpotRoblox(c http.Client) Result {
	resp1, err := GET(c, "https://mathsspot.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	bodyString1 := string(body1)
	if err != nil {
		return Result{Status: StatusFailed}
	}

	if strings.Contains(bodyString1, "FailureServiceNotInRegion") {
		return Result{Status: StatusNo}
	}

	apiPath := extractMathsSpotRobloxAPIPath(bodyString1)
	region := extractMathsSpotRobloxRegion(bodyString1)
	nggFeVersion := extractMathsSpotRobloxNggFeVersion(bodyString1)
	if nggFeVersion == "berlin-v1.34.800_redisexp-arm.1" {
		nggFeVersion = "berlin-v1.34.810_redisexp-arm.1"
	}

	if apiPath == "" || region == "" || nggFeVersion == "" {
		return Result{Status: StatusFailed}
	}

	fakeUAId := genRandomStr(21)
	fakeSessId := genRandomStr(21)
	fakeFesessId := genRandomStr(21)
	fakeVisitId := genRandomStr(21)

	resp2, err := GET(c, "https://mathsspot.com"+apiPath+"/startSession?appId=5349&uaId=ua-"+fakeUAId+"&uaSessionId=uasess-"+fakeSessId+"&feSessionId=fesess-"+fakeFesessId+"&visitId=visitid-"+fakeVisitId+"&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%%3A%%2F%%2Fmathsspot.com%%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA",
		H{"referer", "https://mathsspot.com/"},
		H{"x-ngg-skip-evar-check", "true"},
		H{"x-ngg-fe-version", nggFeVersion},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body2, &res); err != nil {
		nggFeVersion = string(body2)
		resp3, err := GET(c, "https://mathsspot.com"+apiPath+"/startSession?appId=5349&uaId=ua-"+fakeUAId+"&uaSessionId=uasess-"+fakeSessId+"&feSessionId=fesess-"+fakeFesessId+"&visitId=visitid-"+fakeVisitId+"&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%%3A%%2F%%2Fmathsspot.com%%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA",
			H{"referer", "https://mathsspot.com/"},
			H{"x-ngg-skip-evar-check", "true"},
			H{"x-ngg-fe-version", nggFeVersion},
		)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()

		body3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		if err := json.Unmarshal(body3, &res); err != nil {
			return Result{Status: StatusFailed, Err: err}
		}
	}

	switch res.Status {
	case "FailureServiceNotInRegion":
		return Result{Status: StatusNo}
	case "FailureProxyUserLimitExceeded":
		return Result{Status: StatusNo, Info: "Proxy Detected"}
	case "Success", "FailureUnauthorized":
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}

	return Result{Status: StatusFailed}
}
