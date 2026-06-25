package mediaunlocktest

import (
	"net/http"
)

func DSTV(c http.Client) Result {
	return CheckGETStatus(c, "https://now.dstv.com/", ResultMap{
		http.StatusUnavailableForLegalReasons: {Status: StatusNo},
		http.StatusOK:                         {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
