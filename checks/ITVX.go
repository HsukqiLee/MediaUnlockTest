package mediaunlocktest

import (
	"net/http"
)

func ITVX(c http.Client) Result {
	return CheckGETStatus(c, "https://simulcast.itv.com/playlist/itvonline/ITV", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusNotFound:  {Status: StatusOK},
	}, Result{Status: StatusUnexpected}, H{"connection", "keep-alive"})
}
