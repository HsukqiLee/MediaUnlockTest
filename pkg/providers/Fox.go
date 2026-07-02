package providers

import (
	"MediaUnlockTest/pkg/core"
)

func Fox(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://x-live-fox-stgec.uplynk.com/ausw/slices/8d1/d8e6eec26bf544f084bad49a7fa2eac5/8d1de292bcc943a6b886d029e6c0dc87/G00000000.ts?pbs=c61e60ee63ce43359679fb9f65d21564&cloud=aws&si=0", core.ResultMap{
		200: {Status: core.StatusOK},
		403: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
