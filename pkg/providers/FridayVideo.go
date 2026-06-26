package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"errors"
	"io"
)

func FridayVideo(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://video.friday.tw/api2/streaming/get?streamingId=122581&streamingType=2&contentType=4&contentId=1&clientId=",
		core.H{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"},
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
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	if res.Code != "null" {
		switch res.Code {
		case "1006":
			return core.Result{Status: core.StatusNo}
		case "0000":
			return core.Result{Status: core.StatusOK}
		default:
			return core.Result{Status: core.StatusErr, Err: errors.New("unknown code: " + res.Code)}
		}
	}

	return core.Result{Status: core.StatusUnexpected}
}
