package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func Crunchyroll(c core.HttpClient) core.Result {
	resp, err := core.PostForm(c, "https://www.crunchyroll.com/auth/v1/token", `grant_type=client_id`,
		core.H{"Authorization", "Basic Y3Jfd2ViOg=="},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, `"country":"US"`) {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
