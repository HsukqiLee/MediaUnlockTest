package mediaunlocktest

import (
	"net/http"
	"net/http/cookiejar"
)

func DirectvStream(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	return CheckGETStatus(c, "https://stream.directv.com/watchnow", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
