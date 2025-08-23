package mediaunlocktest

import "net/http"

func SHOWTIME(c http.Client) Result {
	resp, err := GET(c, "https://www.paramountpluswithshowtime.com/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
