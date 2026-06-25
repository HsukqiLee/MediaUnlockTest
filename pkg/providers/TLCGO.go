package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"net/http"
	"strings"
)

func TlcGo(c http.Client) core.Result {
	resp, err := core.GET(c, "https://atlas.ngtv.io/v2/locate",
		core.H{"app-id", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuZXR3b3JrIjoiYWxsIiwicHJvZHVjdCI6InByaXNtIiwicGxhdGZvcm0iOiJ3ZWIiLCJhcHBJZCI6ImFsbC1wcmlzbS13ZWItNzI4aGtyIn0.4Fk4E28ffoFgCIcgNSG8xX5TP2n3PIU6c3jadumKULo"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 200 {
		if strings.Contains(string(b), `"country": "US"`) {
			return core.Result{Status: core.StatusOK}
		}
		return core.Result{Status: core.StatusNo}
	}

	return core.Result{Status: core.StatusUnexpected}
}

