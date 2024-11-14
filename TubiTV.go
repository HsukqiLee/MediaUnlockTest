package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func TubiTV(c http.Client) Result {
	resp, err := GET(c, "https://tubitv.com/home")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		resp2, err := GET(c, "https://gdpr.tubi.tv")
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		b, err := io.ReadAll(resp2.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(b), "Unfortunately") {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusOK}
	}
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
