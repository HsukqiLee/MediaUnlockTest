package providers

import (
	"MediaUnlockTest/pkg/core"
)

func J_COM_ON_DEMAND(c core.HttpClient) core.Result {
	c.SetFollowRedirect(true)
	return core.CheckGETStatus(c, "https://linkvod.myjcom.jp/auth/login", core.ResultMap{
		403: {Status: core.StatusNo},
		502: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusOK})
}
