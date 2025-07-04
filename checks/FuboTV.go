package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
	"strconv"
	"math/rand"
)

func FuboTV(c http.Client) Result {
    randNum := strconv.Itoa(rand.Intn(2))
	resp, err := GET(c, "https://api.fubo.tv/appconfig/v1/homepage?platform=web&client_version=R20230310."+randNum+"&nav=v0")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "Forbidden IP") {
		return Result{Status: StatusNo, Info: "IP Forbidden"}
	}
	if strings.Contains(s, "No Subscription") {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}
