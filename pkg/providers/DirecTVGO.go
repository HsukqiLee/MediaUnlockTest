package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func DirecTVGO(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://www.directvgo.com/registrarse", core.ResultMap{
		http.StatusForbidden:        {Status: core.StatusNo},
		http.StatusOK:               {Status: core.StatusNo},
		http.StatusMovedPermanently: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
