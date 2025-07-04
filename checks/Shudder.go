package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func Shudder(c http.Client) Result {
	resp, err := GET(c, "https://www.shudder.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "not available") {
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusOK}
}
