package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func Afreeca(c http.Client) core.Result {
	resp, err := core.GET(c, "https://vod.sooplive.co.kr/player/97464151")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if strings.Contains(bodyString, "document.location.href='https://vod.afreecatv.com'") {
		return core.Result{Status: core.StatusNo}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

