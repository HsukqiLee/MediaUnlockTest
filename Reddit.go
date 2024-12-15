package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func Reddit(c http.Client) Result {
	resp, err := GET(c, "https://www.reddit.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	if resp.StatusCode == 429 {
		return Result{Status: StatusBanned}
	}

	if resp.StatusCode == 200 || resp.StatusCode == 302 {
		return Result{Status: StatusOK}
	}

	if resp.StatusCode == 403 && strings.Contains(bodyString, "blocked") {
		return Result{Status: StatusNo}
	}

	return Result{Status: StatusUnexpected}
}
