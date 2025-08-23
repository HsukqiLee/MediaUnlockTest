package mediaunlocktest

import (
	"net/http"
)

func PandaTV(c http.Client) Result {
	resp, err := GET(c, "https://api.pandalive.co.kr/v1/live/play")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusBadRequest: {Status: StatusOK},
		http.StatusForbidden:  {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
