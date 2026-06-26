package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func SpotvNow(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/5764318566001/videos/6349973203112",
		core.H{"accept", "application/json;pk=BCpkADawqM0U3mi_PT566m5lvtapzMq3Uy7ICGGjGB6v4Ske7ZX_ynzj8ePedQJhH36nym_5mbvSYeyyHOOdUsZovyg2XlhV6rRspyYPw_USVNLaR0fB_AAL2HSQlfuetIPiEzbUs1tpNF9NtQxt3BAPvXdOAsvy1ltLPWMVzJHiw9slpLRgI2NUufc"},
		core.H{"accept-language", "en,zh-CN;q=0.9,zh;q=0.8"},
		core.H{"origin", "https://www.spotvnow.co.kr"},
		core.H{"referer", "https://www.spotvnow.co.kr/"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if strings.Contains(string(b), "CLIENT_GEO") || resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	if resp.StatusCode == 200 || resp.StatusCode == 404 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
