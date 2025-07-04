package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func TubiTV(c http.Client) Result {
	resp1, err := GET(c, "https://tubitv.com")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if resp1.StatusCode == 503 {
		b1, err := io.ReadAll(resp1.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(b1), "geoblock") {
			return Result{Status: StatusNo}
		}
	}
	if resp1.StatusCode == 302 {
		resp2, err := GET(c, "https://gdpr.tubi.tv")
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		b2, err := io.ReadAll(resp2.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(b2), "Unfortunately") {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusOK}
	}
	if resp1.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	if resp1.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
