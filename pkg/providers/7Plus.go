package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func SevenPlus(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://7plus.com.au/", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

