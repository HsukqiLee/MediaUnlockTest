package providers

import (
	"MediaUnlockTest/pkg/core"
	tls_client "github.com/bogdanfinn/tls-client"

	http "github.com/bogdanfinn/fhttp"
)

func DirectvStream(c core.HttpClient) core.Result {
	jar := tls_client.NewCookieJar()
	c.SetCookieJar(jar)
	return core.CheckGETStatus(c, "https://stream.directv.com/watchnow", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
