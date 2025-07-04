package mediaunlocktest

import (
	"net/http"
)

func RakutenMagazine(c http.Client) Result {
	resp, err := GET(c, "https://data-cloudauthoring.magazine.rakuten.co.jp/rem_repository/////////.key")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 403 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 404 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}