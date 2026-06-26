package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

// Project Sekai: Colorful Stage
func PJSK(c core.HttpClient) core.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://game-version.sekai.colorfulpalette.org/1.8.1/3ed70b6a-8352-4532-b819-108837926ff5", nil)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	req.Header.Set("User-Agent", "pjsekai/48 CFNetwork/1240.0.4 Darwin/20.6.0")

	resp, err := core.Cdo(c, req)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return core.Result{Status: core.StatusOK}
	case 403:
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}
