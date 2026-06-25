package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"strings"
)

func ViuCom(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.viu.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if location := resp.Header.Get("location"); location != "" {
		region := strings.Split(location, "/")[4]
		if region == "no-service" {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusOK, Region: region}
	}
	return core.Result{Status: core.StatusNo}
}

