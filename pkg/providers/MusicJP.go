package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
)

func MusicJP(c http.Client) core.Result {
	resp, err := core.GET(c, "https://overseaauth.music-book.jp/globalIpcheck.js")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if string(b) == "" {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}

