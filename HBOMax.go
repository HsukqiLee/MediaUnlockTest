package mediaunlocktest

import (
	"net/http"
	"strings"
)

func HBOMax(c http.Client) Result {
	resp, err := GET(c, "https://www.hbomax.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if strings.Contains(resp.Header.Get("location"), "geo-availability") {
		return Result{Status: StatusNo}
	}
	t := strings.Split(resp.Header.Get("location"), "/")
	region := ""
	if len(t) >= 4 {
		region = strings.Split(resp.Header.Get("location"), "/")[3]
	}
	return Result{Status: StatusOK, Region: strings.ToLower(region)}
}
