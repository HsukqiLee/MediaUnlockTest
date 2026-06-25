package mediaunlocktest

import (
	"net/http"
)

func Channel4(c http.Client) Result {
	return CheckGETStatus(c, "https://www.channel4.com/simulcast/channels/C4", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
