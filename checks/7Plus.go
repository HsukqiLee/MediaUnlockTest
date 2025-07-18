package mediaunlocktest

import (
	"net/http"
)

func SevenPlus(c http.Client) Result {
	return CheckGETStatus(c, "https://7plus.com.au/", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
