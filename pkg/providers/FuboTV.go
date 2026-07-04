package providers

import (
	"MediaUnlockTest/pkg/core"
)

func FuboTV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://api.fubo.tv/popular/unauth/v1?contentType=series&limit=10&genreIds=11")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(
		resp.StatusCode,
		core.ResultMap{
			200: core.Result{Status: core.StatusOK},
			451: core.Result{Status: core.StatusNo},
		},
		core.Result{Status: core.StatusUnexpected},
	)
}
