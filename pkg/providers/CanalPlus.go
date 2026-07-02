package providers

import (
	"MediaUnlockTest/pkg/core"
)

func CanalPlus(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://boutique-tunnel.canalplus.com/", core.ResultMap{
		200: {Status: core.StatusOK},
		302: {Status: core.StatusNo},
		403: {Status: core.StatusBanned},
	}, core.Result{Status: core.StatusUnexpected})
}
