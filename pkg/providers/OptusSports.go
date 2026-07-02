package providers

import (
	"MediaUnlockTest/pkg/core"
)

func OptusSports(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://sport.optus.com.au/api/userauth/validate/web/username/restriction.check@gmail.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		200: {Status: core.StatusOK},
		403: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
