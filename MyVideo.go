package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func MyVideo(c http.Client) Result {
	c.CheckRedirect = nil
	resp, err := GET(c, "https://www.myvideo.net.tw/goLoginPage.do",
		H{"Cookie", "JSESSIONID=12805776F340391EA3CB6A38E12F423A"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if strings.Contains(string(b), "serviceAreaBlock") {
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusOK}
}
