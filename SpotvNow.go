package mediaunlocktest

import (
    "net/http"
    "strings"
    "io"
)

func SpotvNow(c http.Client) Result {
	resp, err := GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/5764318566001/videos/6349973203112",
	    H{"accept", "application/json;pk=BCpkADawqM0U3mi_PT566m5lvtapzMq3Uy7ICGGjGB6v4Ske7ZX_ynzj8ePedQJhH36nym_5mbvSYeyyHOOdUsZovyg2XlhV6rRspyYPw_USVNLaR0fB_AAL2HSQlfuetIPiEzbUs1tpNF9NtQxt3BAPvXdOAsvy1ltLPWMVzJHiw9slpLRgI2NUufc"},
	    H{"accept-language", "en,zh-CN;q=0.9,zh;q=0.8"},
	    H{"origin", "https://www.spotvnow.co.kr"},
	    H{"referer", "https://www.spotvnow.co.kr/"},
    )
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

    if strings.Contains(string(b), "CLIENT_GEO") {
        return Result{Status: StatusNo}
    }
    
    if resp.StatusCode == 200 {
        return Result{Status: StatusOK}
    }

	return Result{Status: StatusUnexpected}
}