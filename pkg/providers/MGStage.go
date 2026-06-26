package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func MGStage(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://www.mgstage.com/", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
