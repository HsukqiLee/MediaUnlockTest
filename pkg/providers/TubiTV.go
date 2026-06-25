package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func TubiTV(c http.Client) core.Result {
	resp1, err := core.GET(c, "https://tubitv.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp1.StatusCode == 503 {
		b1, err := io.ReadAll(resp1.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(b1), "geoblock") {
			return core.Result{Status: core.StatusNo}
		}
	}
	if resp1.StatusCode == 302 {
		resp2, err := core.GET(c, "https://gdpr.tubi.tv")
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		b2, err := io.ReadAll(resp2.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(b2), "Unfortunately") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusOK}
	}
	if resp1.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	if resp1.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}

