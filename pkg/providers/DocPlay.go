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

	if strings.Contains(string(b), "DocPlay hasn't launched in your part of the world yet.") {
		return core.Result{Status: core.StatusNo}
	}

	switch resp.StatusCode {
	case http.StatusTemporaryRedirect:
		if strings.Contains(resp.Header.Get("Location"), "geoblocked") {
			return core.Result{Status: core.StatusNo}
		}
		fallthrough
	case http.StatusSeeOther, http.StatusOK:
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
