package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func CW_TV(c http.Client) Result {
	resp, err := GET(c, "https://www.cwtv.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "302 Found") {
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusOK}
}
