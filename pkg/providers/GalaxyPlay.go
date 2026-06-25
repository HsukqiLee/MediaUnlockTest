package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func GalaxyPlay(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api.glxplay.io/account/device/new")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return core.Result{Status: core.StatusOK}
	}

	if resp.StatusCode == 400 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		var res struct {
			ErrorCode int `json:"errorCode"`
		}
		if len(body) > 0 && body[0] == '<' {
			return core.Result{Status: core.StatusNo}
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return core.Result{Status: core.StatusErr, Err: err}
		}
		if res.ErrorCode == 495 {
			return core.Result{Status: core.StatusNo}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}

