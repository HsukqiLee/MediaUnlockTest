package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func J_COM_ON_DEMAND(c core.HttpClient) core.Result {
	c.SetFollowRedirect(true)
	return core.CheckGETStatus(c, "https://linkvod.myjcom.jp/auth/login", core.ResultMap{
		http.StatusForbidden:  {Status: core.StatusNo},
		http.StatusBadGateway: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusOK})
}
