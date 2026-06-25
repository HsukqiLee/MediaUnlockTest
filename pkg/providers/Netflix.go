package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func NetflixRegion(c http.Client) core.Result {
	// 70143836 绝命毒师
	// 80018499 test
	// 81280792 乐高
	resp, err := core.GET(c, "https://www.netflix.com/title/81280792")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	resp2, err := core.GET(c, "https://www.netflix.com/title/70143836")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	_, err = io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 404 && resp2.StatusCode == 404 {
		return core.Result{Status: core.StatusRestricted, Info: "Originals Only"}
	}
	if resp.StatusCode == 403 && resp2.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}
	if (resp.StatusCode == 200 || resp.StatusCode == 301) || (resp2.StatusCode == 200 || resp2.StatusCode == 301) {
		resp3, err := core.GET(c, "https://www.netflix.com/title/80018499")
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()
		_, err = io.ReadAll(resp3.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		u := resp3.Header.Get("location")
		if u == "" {
			return core.Result{Status: core.StatusOK, Region: "us"}
		}
		t := strings.SplitN(u, "/", 5)
		if len(t) < 5 {
			return core.Result{Status: core.StatusUnexpected}
		}
		return core.Result{Status: core.StatusOK, Region: strings.SplitN(t[3], "-", 2)[0]}
	}
	return core.Result{Status: core.StatusUnexpected}
}

func NetflixCDN(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api.fast.com/netflix/speedtest/v2?https=true&token=YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm&urlCount=5")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned, Info: "IP Banned By Netflix"}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Targets []struct {
			Name     string `json:"name"`
			Url      string `json:"url"`
			Location struct {
				City    string `json:"city"`
				Country string `json:"country"`
			} `json:"location"`
		} `json:"targets"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}
	return core.Result{Status: core.StatusOK, Region: res.Targets[0].Location.Country}
}

