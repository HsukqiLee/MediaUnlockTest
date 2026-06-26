package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func NowTV(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://webtvapi.nowe.com/16/1/getVodURL",
		`{"contentId":"202403181904703","contentType":"Vod","pin":"","deviceName":"Browser","deviceId":"w-663bcc51-913c-913c-913c-913c913c","deviceType":"WEB","secureCookie":null,"callerReferenceNo":"W17151951620081575","profileId":null,"mupId":null}`,
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
		return core.Result{Status: core.StatusUnexpected, Err: err}
	}
	switch res.ResponseCode {
	case "SUCCESS", "ASSET_MISSING", "NOT_LOGIN":
		return core.Result{Status: core.StatusOK}
	case "GEO_CHECK_FAIL":
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}

type noweRes struct {
	ResponseCode string
}
