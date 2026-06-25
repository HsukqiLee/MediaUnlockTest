package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func HoyTV(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://hoytv-live-stream.hoy.tv/ch78/index-fhd.m3u8", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

