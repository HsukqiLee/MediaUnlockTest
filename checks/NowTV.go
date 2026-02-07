package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func NowTV(c http.Client) Result {
	resp, err := PostJson(c, "https://webtvapi.nowe.com/16/1/getVodURL",
		`{"contentId":"202403181904703","contentType":"Vod","pin":"","deviceName":"Browser","deviceId":"w-663bcc51-913c-913c-913c-913c913c","deviceType":"WEB","secureCookie":null,"callerReferenceNo":"W17151951620081575","profileId":null,"mupId":null}`,
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res noweRes
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusUnexpected, Err: err}
	}
	switch res.ResponseCode {
	case "SUCCESS", "ASSET_MISSING", "NOT_LOGIN":
		return Result{Status: StatusOK}
	case "GEO_CHECK_FAIL":
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusUnexpected}
}

type noweRes struct {
	ResponseCode string
}
