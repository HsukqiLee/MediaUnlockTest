package mediaunlocktest

import (
	"net/http"
)

func OptusSports(c http.Client) Result {
	resp, err := GET(c, "https://sport.optus.com.au/api/userauth/validate/web/username/restriction.check@gmail.com")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
