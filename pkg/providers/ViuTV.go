package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func ViuTV(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://api.viu.now.com/p8/3/getLiveURL",
		`{"callerReferenceNo":"20210726112323","contentId":"099","contentType":"Channel","channelno":"099","mode":"prod","deviceId":"29b3cb117a635d5b56","deviceType":"ANDROID_WEB"}`,
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res noweRes
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	switch res.ResponseCode {
	case "SUCCESS":
		return core.Result{Status: core.StatusOK}
	case "GEO_CHECK_FAIL":
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}
