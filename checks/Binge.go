package mediaunlocktest

import (
	"net/http"
)

func Binge(c http.Client) Result {
	return CheckGETStatus(c, "https://auth.streamotion.com.au", ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
