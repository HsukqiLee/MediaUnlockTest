package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func Telasa(c http.Client) core.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api-videopass-anon.kddi-video.com/v1/playback/system_status", nil)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	req.Header.Set("X-Device-ID", "d36f8e6b-e344-4f5e-9a55-90aeb3403799")

	resp, err := core.Cdo(c, req)
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

