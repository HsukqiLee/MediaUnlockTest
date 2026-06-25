package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func AcornTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://acorn.tv/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "Not yet available in your country") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}

