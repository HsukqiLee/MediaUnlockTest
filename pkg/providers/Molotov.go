package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
)

func Molotov(c http.Client) core.Result {
	resp, err := core.GET(c, "https://fapi.molotov.tv/v1/open-europe/is-france")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return core.Result{Status: core.StatusNo}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		IsFrance bool `json:"is_france"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	if res.IsFrance {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusNo}
}

