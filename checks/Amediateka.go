package mediaunlocktest

import (
	"net/http"
)

func Amediateka(c http.Client) Result {
	resp, err := GET(c, "https://www.amediateka.ru/")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 301:
		if resp.Header.Get("Location") == "https://www.amediateka.ru/unavailable/index.html?page=https://www.amediateka.ru/" {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusUnexpected}
	case 200:
		return Result{Status: StatusOK}
	case 503, 445:
		return Result{Status: StatusBanned}
	default:
		return Result{Status: StatusUnexpected}
	}
}
