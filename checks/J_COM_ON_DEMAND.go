package mediaunlocktest

import (
	"net/http"
)

func J_COM_ON_DEMAND(c http.Client) Result {
	c.CheckRedirect = nil
	return CheckGETStatus(c, "https://linkvod.myjcom.jp/auth/login", ResultMap{
		http.StatusForbidden:  {Status: StatusNo},
		http.StatusBadGateway: {Status: StatusNo},
	}, Result{Status: StatusOK})
}
