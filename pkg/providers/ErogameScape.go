package providers

import (
	"MediaUnlockTest/pkg/core"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func ErogameScape(c core.HttpClient) core.Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://erogamescape.org/~ap2/ero/toukei_kaiseki/", nil)
	resp, err := c.Do(req)

	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, context.DeadlineExceeded) {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "18歳") {
			return core.Result{Status: core.StatusOK}
		}
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}
