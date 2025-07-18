package mediaunlocktest

import (
	"net/http"
)

func JioCinema(c http.Client) Result {
	return CheckGETStatus(c, "https://content-jiovoot.voot.com/psapi/", ResultMap{
		http.StatusOK: {Status: StatusOK},
		474:           {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
