package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func FOD(c http.Client) core.Result {
	resp, err := core.GET(c, "https://geocontrol1.stream.ne.jp/fod-geo/check.xml?time=1624504256")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "true") {
		return core.Result{Status: core.StatusOK}
	}
	if strings.Contains(s, "false") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}

