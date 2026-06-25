package mediaunlocktest

import (
	"net/http"
)

func KOCOWA(c http.Client) Result {
	return CheckGETStatus(c, "https://www.kocowa.com/", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
