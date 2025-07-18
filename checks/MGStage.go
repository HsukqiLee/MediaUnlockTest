package mediaunlocktest

import (
	"net/http"
)

func MGStage(c http.Client) Result {
	return CheckGETStatus(c, "https://www.mgstage.com/", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
