package providers

import (
	"MediaUnlockTest/pkg/core"
	"strings"
)

func Hotstar(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.hotstar.com/o/v1/page/1557?offset=0&size=20&tao=0&tas=20")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 475:
		return core.Result{Status: core.StatusNo}
	case 401:
		resp, err := core.GET(c, "https://www.hotstar.com")
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		if resp.StatusCode == 301 {
			return core.Result{Status: core.StatusNo}
		}
		u := resp.Header.Get("Location")
		if u == "" {
			return core.Result{Status: core.StatusNo}
		}
		t := strings.SplitN(u, "/", 4)
		if len(t) < 4 {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusOK, Region: t[3]}
	case 472, 473, 474:
		return core.Result{Status: core.StatusBanned}
	}
	return core.Result{Status: core.StatusUnexpected}
}
