package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func JioCinema(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://content-jiovoot.voot.com/psapi/", core.ResultMap{
		http.StatusOK: {Status: core.StatusOK},
		474:           {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
