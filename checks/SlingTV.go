package mediaunlocktest

import "net/http"

func SlingTV(c http.Client) Result {
	resp, err := GET(c, "https://www.sling.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
		http.StatusFound:     {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
