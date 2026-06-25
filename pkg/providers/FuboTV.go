package providers

import (
	"MediaUnlockTest/pkg/core"
	"crypto/rand"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func FuboTV(c http.Client) core.Result {
	// Generate cryptographically secure random number
	randomBytes := make([]byte, 1)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	randNum := strconv.Itoa(int(randomBytes[0]) % 2)
	resp, err := core.GET(c, "https://api.fubo.tv/appconfig/v1/homepage?platform=web&client_version=R20230310."+randNum+"&nav=v0")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "Forbidden IP") {
		return core.Result{Status: core.StatusNo, Info: "IP Forbidden"}
	}
	if strings.Contains(s, "No Subscription") {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}

