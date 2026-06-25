package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func MyVideo(c http.Client) core.Result {
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := core.GET(c, "https://www.myvideo.net.tw/login.do")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		if resp.Header.Get("Location") == "/serviceAreaBlock.do" {
			return core.Result{Status: core.StatusNo}
		}
		if resp.Header.Get("Location") == "/goLoginPage.do" {
			return core.Result{Status: core.StatusOK}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}

