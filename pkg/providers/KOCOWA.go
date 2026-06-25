package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func KOCOWA(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://www.kocowa.com/", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

