package mediaunlocktest

import (
    "net/http"
    "strings"
    "io"
)

func MaoriTV(c http.Client) Result {
    resp, err := GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/1614493167001/videos/6278939271001",
        H{"Accept", "application/json;pk=BCpkADawqM2E9yW4lLgKIEIV5majz5djzZCIqJiYMkP5yYaYdF6AQYq4isPId1ZLtQdGnK1ErLYG0-r1N-3DzAEdbfvw9SFdDWz_i09pLp8Njx1ybslyIXid-X_Dx31b7-PLdQhJCws-vk6Y"},
        H{"Origin", "https://www.maoritelevision.com"},
    )
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    
    if strings.Contains(string(b), "CLIENT_GEO") {
        return Result{Status: StatusNo}
    }
    
    if resp.StatusCode == 200 {
        return Result{Status: StatusOK}
    }
    
    return Result{Status: StatusUnexpected}
}
