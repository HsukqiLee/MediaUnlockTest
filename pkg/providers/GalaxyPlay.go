package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func GalaxyPlay(c http.Client) Result {
	resp, err := GET(c, "https://api.glxplay.io/account/device/new")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return Result{Status: StatusOK}
	}

	if resp.StatusCode == 400 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		var res struct {
			ErrorCode int `json:"errorCode"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			if err.Error() == "invalid character 'S' looking for beginning of value" {
				return Result{Status: StatusNo}
			}
			return Result{Status: StatusErr, Err: err}
		}
		if res.ErrorCode == 495 {
			return Result{Status: StatusNo}
		}
	}
	return Result{Status: StatusUnexpected}
}
