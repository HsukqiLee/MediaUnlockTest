package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Channel4(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://www.channel4.com/simulcast/channels/C4", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

