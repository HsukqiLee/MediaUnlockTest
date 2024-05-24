package mediaunlocktest

import (
	"net/http"
)

func OptusSports(c http.Client) Result {
	resp, err := GET(c, "https://sport.optus.com.au/api/userauth/validate/web/username/restriction.check@gmail.com")
	if err != nil {
		return Result{Status: StatusNetworkErr}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}