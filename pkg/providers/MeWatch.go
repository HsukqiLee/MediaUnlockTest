package providers

import (
	"MediaUnlockTest/pkg/core"
)

func MeWatch(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://cdn.mewatch.sg/api/items/97098/videos?delivery=stream%2Cprogressive&ff=idp%2Cldp%2Crpt%2Ccd&lang=en&resolution=External&segments=all", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
