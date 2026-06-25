package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func Radiko(c http.Client) core.Result {
	resp, err := core.GET(c, "https://radiko.jp/area?_=1625406539531")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, `classs="OUT"`) {
		return core.Result{Status: core.StatusNo}
	}
	if strings.Contains(s, "JAPAN") {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}

