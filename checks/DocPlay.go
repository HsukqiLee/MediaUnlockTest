package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func DocPlay(c http.Client) Result {
	resp, err := GET(c, "https://www.docplay.com/subscribe")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	body := string(b)

	if strings.Contains(body, "DocPlay hasn't launched in your part of the world yet.") {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusOK {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
