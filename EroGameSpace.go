package mediaunlocktest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func EroGameSpace(c http.Client) Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://erogamescape.org/~ap2/ero/toukei_kaiseki/", nil)
	resp, err := c.Do(req)

	if err != nil {
		if err.Error() == `Get "https://erogamescape.org/~ap2/ero/toukei_kaiseki/": context deadline exceeded` || err.Error() == `Get "https://erogamescape.org/~ap2/ero/toukei_kaiseki/": EOF` {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "18歳") {
			return Result{Status: StatusOK}
		}
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
