package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func U_NEXT(c http.Client) Result {
	resp, err := PostJson(c, "https://cc.unext.jp", `{"operationName":"cosmo_getPlaylistUrl","variables":{"code":"ED00479780","playMode":"caption","bitrateLow":192,"bitrateHigh":null,"validationOnly":false},"query":"query cosmo_getPlaylistUrl($code: String, $playMode: String, $bitrateLow: Int, $bitrateHigh: Int, $validationOnly: Boolean) {\n  webfront_playlistUrl(\n    code: $code\n    playMode: $playMode\n    bitrateLow: $bitrateLow\n    bitrateHigh: $bitrateHigh\n    validationOnly: $validationOnly\n  ) {\n    subTitle\n    playToken\n    playTokenHash\n    beaconSpan\n    result {\n      errorCode\n      errorMessage\n      __typename\n    }\n    resultStatus\n    licenseExpireDate\n    urlInfo {\n      code\n      startPoint\n      resumePoint\n      endPoint\n      endrollStartPosition\n      holderId\n      saleTypeCode\n      sceneSearchList {\n        IMS_AD1\n        IMS_L\n        IMS_M\n        IMS_S\n        __typename\n      }\n      movieProfile {\n        cdnId\n        type\n        playlistUrl\n        movieAudioList {\n          audioType\n          __typename\n        }\n        licenseUrlList {\n          type\n          licenseUrl\n          __typename\n        }\n        __typename\n      }\n      umcContentId\n      movieSecurityLevelCode\n      captionFlg\n      dubFlg\n      commodityCode\n      movieAudioList {\n        audioType\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"}`)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	
	var res struct {
		Data struct {
			PlaylistUrl struct {
				ResultStatus int `json:"resultStatus"`
			} `json:"webfront_playlistUrl"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}
	
	status := res.Data.PlaylistUrl.ResultStatus
	if status == 200 || status == 475 {
		return Result{Status: StatusOK}
	}
	
	if status == 467{
		return Result{Status: StatusNo}
	}
	
	return Result{Status: StatusUnexpected}
}

