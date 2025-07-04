package mediaunlocktest

import (
	"net/http"
	"net/http/cookiejar"
)

func DirectvStream(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	resp, err := GET(c, "https://stream.directv.com/watchnow")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
