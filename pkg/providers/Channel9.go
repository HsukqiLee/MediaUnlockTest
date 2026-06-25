package mediaunlocktest

import (
	"net/http"
)

func Channel9(c http.Client) Result {
	return CheckGETStatus(c, "https://login.nine.com.au", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusFound:     {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
