package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func HamiVideo(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://hamivideo.hinet.net/api/play.do?id=OTT_VOD_0000249064&freeProduct=1")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Code string
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Code == "06001-107" {
		return core.Result{Status: core.StatusOK}
	}
	if res.Code == "06001-106" {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}
