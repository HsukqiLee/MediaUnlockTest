package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func Philo(c http.Client) core.Result {
	resp, err := core.GET(c, "https://content-us-east-2-fastly-b.www.philo.com/geo")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if strings.Contains(string(b), "SUCCESS") {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}

