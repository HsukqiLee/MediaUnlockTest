package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	//"log"
)

func NowE(c http.Client) Result {
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
	//log.Println(string(b))
	var res noweRes
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusUnexpected, Err: err}
	}
	if res.ResponseCode == "SUCCESS" {
		return Result{Status: StatusOK}
	} else if res.ResponseCode == "GEO_CHECK_FAIL" {
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusUnexpected}
}

type noweRes struct {
	ResponseCode string
}
