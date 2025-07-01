package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func ClipTV(c http.Client) Result {
	resp, err := GET(c, "https://cliptv.vn/truyen-hinh")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		bodyString := string(body)
		if strings.Contains(bodyString, "Sorry, this video is not available in your country.") {
			return Result{Status: StatusNo}
		}

		return Result{Status: StatusOK}
	}
	return Result{Status: StatusUnexpected}
}
