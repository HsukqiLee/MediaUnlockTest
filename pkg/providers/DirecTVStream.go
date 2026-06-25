package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"net/http/cookiejar"
)

func DirectvStream(c http.Client) core.Result {
	c.Jar, _ = cookiejar.New(nil)
	return core.CheckGETStatus(c, "https://stream.directv.com/watchnow", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

