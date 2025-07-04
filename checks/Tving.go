package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func Tving(c http.Client) Result {
	resp, err := GET(c, "https://api.tving.com/v2a/media/stream/info?apiKey=1e7952d0917d6aab1f0293a063697610&mediaCode=RV60891248")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res struct {
		Body struct {
			Result struct {
				Code string `json:"code"`
			} `json:"result"`
		} `json:"body"`
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	if res.Body.Result.Code == "001" {
		return Result{Status: StatusNo}
	}

	if res.Body.Result.Code == "000" {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
