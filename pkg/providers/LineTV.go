package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func LineTV(c http.Client) core.Result {
	// GET the drama/episode page which has SSR data in window.__INITIAL_STATE__
	resp, err := core.GET(c, "https://www.linetv.tw/drama/11829/eps/1")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString := string(body)

	// The page uses SSR and injects window.__INITIAL_STATE__ with drama episode info.
	// If the IP is outside TW, the page either redirects or returns empty eps_info.
	// Check for the presence of eps_info with episode data.
	if !strings.Contains(bodyString, "window.__INITIAL_STATE__") {
		return core.Result{Status: core.StatusNetworkErr}
	}

	// Check for eps_info with at least one episode entry (indicated by "durationInMs")
	reEps := regexp.MustCompile(`"eps_info"\s*:\s*\[`)
	reMs := regexp.MustCompile(`"durationInMs"\s*:\s*\d+`)
	if reEps.MatchString(bodyString) && reMs.MatchString(bodyString) {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusNo}
}
