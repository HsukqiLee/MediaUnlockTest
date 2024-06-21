package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Epix(c http.Client) Result {
	url := "https://api.epix.com/v2/sessions"
	resp, err := PostJson(c, url,
		`{"device":{"guid":"7a0baaaf-384c-45cd-a21d-310ca5d3002a","format":"console","os":"web","display_width":1865,"display_height":942,"app_version":"1.0.2","model":"browser","manufacturer":"google"},"apikey":"53e208a9bbaee479903f43b39d7301f7"}`,
		H{"connection", "keep-alive"},
		H{"traceparent", "00-000000000000000015b7efdb572b7bf2-4aefaea90903bd1f-01"},
		H{"x-datadog-sampling-priority", "1"},
		H{"x-datadog-trace-id", "1564983120873880562"},
		H{"x-datadog-parent-id", "5399726519264460063"},
		H{"origin", "https://www.mgmplus.com"},
		H{"referer", "https://www.mgmplus.com/"},
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
	if strings.Contains(s, "error code") {
		return Result{Status: StatusNo}
	}
	if strings.Contains(s, "blocked") {
		return Result{Status: StatusBanned}
	}
	var res struct {
		DeviceSession struct {
			SessionToken string `json:"session_token"`
		} `json:"device_session"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		// log.Println(err)
		return Result{Status: StatusFailed, Err: err}
	}
	url2 := "https://api.epix.com/v2/movies/16921/play"
	resp2, err := PostJson(c, url2, `{}`, H{"X-Session-Token", res.DeviceSession.SessionToken})
	if err != nil {
		// log.Println(err)
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		// log.Println(err)
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res2 struct {
		Movie struct {
			Entitlements struct {
				Status string
			}
		}
	}
	if err := json.Unmarshal(b2, &res2); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	switch res2.Movie.Entitlements.Status {
	case "PROXY_DETECTED":
		return Result{Status: StatusNo, Info: "Proxy Detected"}
	case "GEO_BLOCKED":
		return Result{Status: StatusNo}
	case "NOT_SUBSCRIBED":
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusFailed}
}
