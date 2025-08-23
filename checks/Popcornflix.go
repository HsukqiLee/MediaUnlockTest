package mediaunlocktest

import (
	"net/http"
)

func Popcornflix(c http.Client) Result {
	resp, err := GET(c, "https://popcornflix-prod.cloud.seachange.com/cms/popcornflix/clientconfiguration/versions/2")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
