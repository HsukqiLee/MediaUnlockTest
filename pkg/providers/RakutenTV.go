package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func RakutenTV_EU(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://gizmo.rakuten.tv/v3/me/start?device_identifier=web&device_stream_audio_quality=2.0&device_stream_hdr_type=NONE&device_stream_video_quality=FHD", `{"device_identifier":"web","device_metadata":{"app_version":"v5.5.22","audio_quality":"2.0","brand":"chrome","firmware":"XX.XX.XX","hdr":false,"model":"GENERIC","os":"Android OS","sdk":"112.0.0","serial_number":"not implemented","trusted_uid":false,"uid":"ab0dd3e8-5cae-4ad2-ba86-97af867e75c3","video_quality":"FHD","year":1970},"ifa_id":"b9c55e58-d5d0-41ed-becb-a54499731531"}`)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 406 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if strings.Contains(bodyString, "forbidden_market") {
		return core.Result{Status: core.StatusNo}
	}

	if strings.Contains(bodyString, "forbidden_vpn") {
		return core.Result{Status: core.StatusNo, Info: "VPN Forbidden"}
	}

	return core.Result{Status: core.StatusUnexpected}
}

func RakutenTV_JP(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.tv.rakuten.co.jp/content/playinfo.json?content_id=476611&device_id=14&trailer=1&auth=0&log=0&serial_code=&tmp_eng_flag=1&multi_audio_support=1&_=1716694365356",
		core.H{"connection", "keep-alive"},
		core.H{"Cookie", "alt_id=kdPG3ErDszsWchi~f3P7Y3Mk; _ra=1716693934724|fbf06bf6-0e63-49bc-b5ae-ea8e785126ba; sec_token=6d518581124ba17c1b9968dca83aba7d441dcf88s%3A40%3A%220f817994db4925695da3375e3248a7552d981647%22%3B"},
		core.H{"origin", "https://tv.rakuten.co.jp"},
		core.H{"referer", "https://tv.rakuten.co.jp/"},
	)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
