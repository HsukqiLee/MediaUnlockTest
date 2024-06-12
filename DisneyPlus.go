package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func extractDisneyPlusRegion(body string) string {
    re := regexp.MustCompile(`"countryCode"\s*:\s*"([^"]+)"`)
    match := re.FindStringSubmatch(body)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}

func extractDisneyPlusSupport(body string) bool {
    re := regexp.MustCompile(`"inSupportedLocation"\s*:\s*(false|true)`)
    match := re.FindStringSubmatch(body)
    if len(match) > 1 && match[1] == "true" {
        return true
    }
    return false
}

func DisneyPlus(c http.Client) Result {
	resp1, err := PostJson(c, "https://disney.api.edge.bamgrid.com/devices",
	`{"deviceFamily":"browser","applicationRuntime":"chrome","deviceProfile":"windows","attributes":{}}`,
	     H{"authorization", "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    bodyString1 := string(body1)
	if strings.Contains(bodyString1, "403 ERROR") {
		return Result{Status: StatusNo}
	}
	
	var res1 struct {
		Assertion string `json:"assertion"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	resp2, err := PostForm(c, "https://disney.api.edge.bamgrid.com/token",
	    `grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Atoken-exchange&latitude=0&longitude=0&platform=browser&subject_token=` + res1.Assertion + `&subject_token_type=urn%3Abamtech%3Aparams%3Aoauth%3Atoken-type%3Adevice`,
	    H{"authorization", "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    
    bodyString2 := string(body2)
	if strings.Contains(bodyString2, "forbidden-location") || resp2.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	resp3, err := GET(c, "https://disneyplus.com")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()
    if strings.Contains(resp3.Request.URL.String(), "preview") || strings.Contains(resp3.Request.URL.String(), "unavailable") {
		return Result{Status: StatusNo}
	}
    
	var res2 struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(body2, &res2); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	resp4, err := PostJson(c, "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
	    `{"query":"mutation refreshToken($input: RefreshTokenInput!) {\n            refreshToken(refreshToken: $input) {\n                activeSession {\n                    sessionId\n                }\n            }\n        }","variables":{"input":{"refreshToken":"` + res2.RefreshToken + `"}}}`,
	    H{"authorization", "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()
	body4, err := io.ReadAll(resp4.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString4 := string(body4)

    if !extractDisneyPlusSupport(bodyString4) {
        return Result{Status: StatusNo}
    }

    region := extractDisneyPlusRegion(bodyString4)
	if region == "" {
		return Result{Status: StatusUnexpected}
	}
    
    return Result{Status: StatusOK, Region: strings.ToLower(region)}
}
