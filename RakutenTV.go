package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func RakutenTV(c http.Client) Result {
    resp, err := PostJson(c, "https://gizmo.rakuten.tv/v3/me/start?device_identifier=web&device_stream_audio_quality=2.0&device_stream_hdr_type=NONE&device_stream_video_quality=FHD", `{"device_identifier":"web","device_metadata":{"app_version":"v5.5.22","audio_quality":"2.0","brand":"chrome","firmware":"XX.XX.XX","hdr":false,"model":"GENERIC","os":"Android OS","sdk":"112.0.0","serial_number":"not implemented","trusted_uid":false,"uid":"ab0dd3e8-5cae-4ad2-ba86-97af867e75c3","video_quality":"FHD","year":1970},"ifa_id":"b9c55e58-d5d0-41ed-becb-a54499731531"}`)
    
    if err != nil {
        return Result{Status: StatusNetworkErr}
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
    	return Result{Status: StatusOK}
	}
    
    bodyBytes, err := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    
    if err != nil {
		return Result{Status: StatusFailed}
	}
	
	if strings.Contains(bodyString, "forbidden_market") {
		return Result{Status: StatusNo, Info: "Not Available"}
	}
	
	if strings.Contains(bodyString, "forbidden_vpn") {
		return Result{Status: StatusNo, Info: "VPN Forbidden"}
	}
	
	return Result{Status: StatusUnexpected}
}