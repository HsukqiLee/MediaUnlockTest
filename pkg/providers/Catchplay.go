package mediaunlocktest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func Catchplay(c http.Client) Result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://sunapi.catchplay.com/geo", nil)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	req.Header.Set("authorization", "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-TW,zh;q=0.9,en;q=0.8")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := cdo(c, req)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		Code string `json:"code"`
		Data struct {
			IsoCode string `json:"isoCode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}
	if res.Code == "100016" {
		return Result{Status: StatusNo}
	}
	region := res.Data.IsoCode
	if region != "" {
		return Result{Status: StatusOK, Region: strings.ToLower(region)}
	}

	return Result{Status: StatusUnexpected}
}
