package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func ZDF(c http.Client) core.Result {
	return core.CheckDalvikStatus(c, "https://ssl.zdf.de/geo/de/geo.txt", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

