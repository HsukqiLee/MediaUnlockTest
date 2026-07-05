package providers

import (
	"MediaUnlockTest/pkg/core"
	"strings"
)

func Grok(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://grok.com/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://grok.com/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	res := core.ResultFromMapping(
		resp.StatusCode,
		core.ResultMap{
			200: core.Result{Status: core.StatusOK},
			403: core.Result{Status: core.StatusNo},
		},
		core.Result{Status: core.StatusUnexpected},
	)
	if res.Status == core.StatusOK || res.Status == core.StatusNo {
		res.Region = strings.ToLower(loc)
	}
	return res
}
