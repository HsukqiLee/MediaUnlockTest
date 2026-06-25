package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Channel9(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://login.nine.com.au", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusFound:     {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

