package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func Molotov(c http.Client) Result {
	resp, err := GET(c, "https://fapi.molotov.tv/v1/open-europe/is-france")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res struct {
		IsFrance bool `json:"is_france"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	if res.IsFrance {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusNo}
}
