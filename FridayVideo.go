package mediaunlocktest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FridayVideo(c http.Client) Result {
	resp, err := GET(c, "https://video.friday.tw/api2/streaming/get?streamingId=122581&streamingType=2&contentType=4&contentId=1&clientId=",
		H{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	if res.Code != "null" {
		switch res.Code {
		case "1006":
			return Result{Status: StatusNo}
		case "0000":
			return Result{Status: StatusOK}
		default:
			return Result{Status: StatusErr, Err: fmt.Errorf("unknown code: %s", res.Code)}
		}
	}

	return Result{Status: StatusUnexpected}
}
