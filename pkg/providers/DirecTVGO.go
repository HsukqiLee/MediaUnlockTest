package providers

import (
	"MediaUnlockTest/pkg/core"
)

func DirecTVGO(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://www.directvgo.com/registrarse", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusNo},
		301: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
