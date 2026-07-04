package providers

import "MediaUnlockTest/pkg/core"

func Perplexity(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.perplexity.ai/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(
		resp.StatusCode,
		core.ResultMap{
			200: core.Result{Status: core.StatusOK},
			403: core.Result{Status: core.StatusNo},
		},
		core.Result{Status: core.StatusUnexpected},
	)
}
