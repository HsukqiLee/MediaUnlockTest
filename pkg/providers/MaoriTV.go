package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func MaoriTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/1614493167001/videos/6278939271001",
		core.H{"Accept", "application/json;pk=BCpkADawqM2E9yW4lLgKIEIV5majz5djzZCIqJiYMkP5yYaYdF6AQYq4isPId1ZLtQdGnK1ErLYG0-r1N-3DzAEdbfvw9SFdDWz_i09pLp8Njx1ybslyIXid-X_Dx31b7-PLdQhJCws-vk6Y"},
		core.H{"Origin", "https://www.maoritelevision.com"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if strings.Contains(string(b), "CLIENT_GEO") {
		return core.Result{Status: core.StatusNo}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

