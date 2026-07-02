package providers

import (
	"MediaUnlockTest/pkg/core"
)

// World Flipper Japan
func WFJP(c core.HttpClient) core.Result {
	return core.CheckDalvikStatus(c, "https://api.worldflipper.jp/", core.ResultMap{
		200: {Status: core.StatusOK},
		403: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
