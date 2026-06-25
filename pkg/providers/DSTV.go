package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func DSTV(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://now.dstv.com/", core.ResultMap{
		http.StatusUnavailableForLegalReasons: {Status: core.StatusNo},
		http.StatusOK:                         {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

