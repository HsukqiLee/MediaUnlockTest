package providers

import "MediaUnlockTest/pkg/core"

func Mistral(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://chat.mistral.ai/")
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
			307: core.Result{Status: core.StatusOK},
			308: core.Result{Status: core.StatusOK},
			403: core.Result{Status: core.StatusNo},
		},
		core.Result{Status: core.StatusUnexpected},
	)
}
