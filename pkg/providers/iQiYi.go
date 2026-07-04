package providers

import (
	"MediaUnlockTest/pkg/core"
	"strings"
)

func IQiYi(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.iq.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	s := resp.Header.Get("x-custom-client-ip")
	if s == "" {
		return core.Result{Status: core.StatusNo}
	}
	_, after, ok := strings.Cut(s, ":")
	if !ok {
		return core.Result{Status: core.StatusNo}
	}
	region := after
	if region == "ntw" {
		region = "tw"
	}
	return core.Result{Status: core.StatusOK, Region: region}
}
