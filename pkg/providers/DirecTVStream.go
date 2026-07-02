package providers

import (
	"MediaUnlockTest/pkg/core"

	tls_client "github.com/bogdanfinn/tls-client"
)

func DirectvStream(c core.HttpClient) core.Result {
	jar := tls_client.NewCookieJar()
	c.SetCookieJar(jar)
	return core.CheckGETStatus(c, "https://stream.directv.com/watchnow", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
