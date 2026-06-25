package providers

import (
	"net/http"
	"MediaUnlockTest/pkg/core"
)

// World Flipper Japan
func WFJP(c http.Client) core.Result {
	return core.CheckDalvikStatus(c, "https://api.worldflipper.jp/", core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

