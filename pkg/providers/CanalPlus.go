package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func CanalPlus(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://boutique-tunnel.canalplus.com/", core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusFound:     {Status: core.StatusNo},
		http.StatusForbidden: {Status: core.StatusBanned},
	}, core.Result{Status: core.StatusUnexpected})
}
