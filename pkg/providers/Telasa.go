package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func Telasa(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api-videopass-anon.kddi-video.com/v1/playback/system_status", core.H{"X-Device-ID", "d36f8e6b-e344-4f5e-9a55-90aeb3403799"})
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Status struct {
			Type    string
			Subtype string
		}
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	if res.Status.Subtype == "IPLocationNotAllowed" {
		return core.Result{Status: core.StatusNo}
	}
	if res.Status.Type != "" {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
