package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func DocPlay(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.docplay.com/subscribe")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	body := string(b)

	if strings.Contains(body, "DocPlay hasn't launched in your part of the world yet.") {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusOK {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
