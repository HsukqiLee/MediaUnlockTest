package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"strings"
)

func Steam(c http.Client) core.Result {
	resp, err := core.GET(c, "https://store.steampowered.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "steamCountry" {
			i := strings.Index(c.Value, "%")
			if i == -1 {
				return core.Result{Status: core.StatusNo}
			}
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(c.Value[:i])}
		}
	}
	return core.Result{Status: core.StatusNo}
}

