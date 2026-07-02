package providers

import (
	"MediaUnlockTest/pkg/core"
)

func HoyTV(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://hoytv-live-stream.hoy.tv/ch78/index-fhd.m3u8", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
