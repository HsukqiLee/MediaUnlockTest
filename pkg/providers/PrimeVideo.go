package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func PrimeVideo(c http.Client) core.Result {
	c.CheckRedirect = nil
	resp, err := core.GET(c, "https://www.primevideo.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if i := strings.Index(s, `"currentTerritory":`); i != -1 {
		return core.Result{
			Status: core.StatusOK,
			Region: strings.ToLower(s[i+20 : i+22]),
		}
	}
	return core.Result{Status: core.StatusNo}
}

