package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func J_COM_ON_DEMAND(c http.Client) core.Result {
	c.CheckRedirect = nil
	return core.CheckGETStatus(c, "https://linkvod.myjcom.jp/auth/login", core.ResultMap{
		http.StatusForbidden:  {Status: core.StatusNo},
		http.StatusBadGateway: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusOK})
}

