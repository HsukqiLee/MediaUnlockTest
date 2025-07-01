package mediaunlocktest

import (
	"net/http"
)

func Sooka(c http.Client) Result {
	resp, err := GET(c, "https://sooka.my/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return Result{Status: StatusOK}
	case 403:
		return Result{Status: StatusNo}
	default:
		return Result{Status: StatusUnexpected}
	}
}
