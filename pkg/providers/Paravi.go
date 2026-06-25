package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func Paravi(c http.Client) core.Result {
	resp, err := core.PostJson(c, "https://api.paravi.jp/api/v1/playback/auth",
		`{"meta_id":17414,"vuid":"3b64a775a4e38d90cc43ea4c7214702b","device_code":1,"app_id":1}`,
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	
	var res struct {
		Error struct {
			Type string
		}
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	
	if res.Error.Type == "Unauthorized" {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}

