package mediaunlocktest

import (
	"net/http"
)

func MyVideo(c http.Client) Result {
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := GET(c, "https://www.myvideo.net.tw/login.do")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		if resp.Header.Get("Location") == "/serviceAreaBlock.do" {
			return Result{Status: StatusNo}
		}
		if resp.Header.Get("Location") == "/goLoginPage.do" {
			return Result{Status: StatusOK}
		}
	}
	return Result{Status: StatusUnexpected}
}
