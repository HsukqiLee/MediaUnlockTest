package mediaunlocktest

import (
	"encoding/json"
    "io"
	"net/http"
	"strings"
)

func Channel10(c http.Client) Result {
	resp, err := GET(c, "https://10play.com.au/geo-web")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		allow bool
	}
	if err := json.Unmarshal(b, &res); err != nil {
	    if strings.Contains(string(b), "not available") {
	        return Result{Status: StatusNo}
        }
		return Result{Status: StatusErr, Err: err}
	}
	if res.allow {
		return Result{Status: StatusOK}
	}
	return Result{Status: StatusNo}
}