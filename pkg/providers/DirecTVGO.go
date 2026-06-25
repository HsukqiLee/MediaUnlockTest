package mediaunlocktest

import (
	"net/http"
)

func DirecTVGO(c http.Client) Result {
	return CheckGETStatus(c, "https://www.directvgo.com/registrarse", ResultMap{
		http.StatusForbidden:        {Status: StatusNo},
		http.StatusOK:               {Status: StatusNo},
		http.StatusMovedPermanently: {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
