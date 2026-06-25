package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func Crunchyroll(c http.Client) Result {
	resp, err := PostForm(c, "https://www.crunchyroll.com/auth/v1/token", `grant_type=client_id`,
		H{"Authorization", "Basic Y3Jfd2ViOg=="},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, `"country":"US"`) {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}
