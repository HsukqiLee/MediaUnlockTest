package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func LiTV(c http.Client) core.Result {
	deviceID, err := getLiTVDeviceID(c)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	payload := fmt.Sprintf(
		`{"jsonrpc":"2.0","id":0,"method":"CCCService.GetProgramInformation","params":{"version":"2.0","project_num":"LTWEB02","device_id":"%s","swver":"LTWEB0210000WEB20190612185813","content_id":"VOD00328856","content_type":"drama"}}`,
		deviceID,
	)

	resp, err := core.PostJson(c, "https://proxy.svc.litv.tv/cdi/v2/rpc", payload,
		core.H{"Origin", "https://www.litv.tv"},
		core.H{"Referer", "https://www.litv.tv/drama/watch/VOD00328856"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Result *struct {
			Data *struct {
				ContentID string `json:"content_id"`
			} `json:"data"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return core.Result{Status: core.StatusUnexpected, Err: err}
	}

	if res.Error != nil {
		// code -32000 or similar = service error, probably geo-blocked
		return core.Result{Status: core.StatusNo}
	}

	if res.Result != nil && res.Result.Data != nil && res.Result.Data.ContentID != "" {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusNo}
}

func getLiTVDeviceID(c http.Client) (string, error) {
	resp, err := core.PostJson(c, "https://www.litv.tv/api/generate-device-id", "",
		core.H{"Origin", "https://www.litv.tv"},
		core.H{"Referer", "https://www.litv.tv/"},
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		DeviceId string `json:"deviceId"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	if res.DeviceId == "" {
		return "", fmt.Errorf("deviceId not found in response")
	}
	return res.DeviceId, nil
}
