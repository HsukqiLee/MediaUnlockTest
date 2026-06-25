package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func Afreeca(c http.Client) Result {
	resp, err := GET(c, "https://vod.sooplive.co.kr/player/97464151")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if strings.Contains(bodyString, "document.location.href='https://vod.afreecatv.com'") {
		return Result{Status: StatusNo}
	}

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK: {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
