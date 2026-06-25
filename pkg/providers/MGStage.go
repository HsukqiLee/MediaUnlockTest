package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func MGStage(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://www.mgstage.com/", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

