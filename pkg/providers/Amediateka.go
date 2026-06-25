package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Amediateka(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.amediateka.ru/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 301, 302:
		if resp.Header.Get("Location") == "https://www.amediateka.ru/unavailable/index.html?page=https://www.amediateka.ru/" {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusUnexpected}
	case 200:
		return core.Result{Status: core.StatusOK}
	case 503, 445:
		return core.Result{Status: core.StatusBanned}
	default:
		return core.Result{Status: core.StatusUnexpected}
	}
}

