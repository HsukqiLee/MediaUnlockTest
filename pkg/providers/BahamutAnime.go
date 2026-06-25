package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func BahamutAnime(c http.Client) core.Result {
	c.Jar, _ = cookiejar.New(nil)

	headers := core.GetRealisticHeaders("html")
	headers = append(headers, core.H{"x-custom-headers", "true"})

	type apiResponse struct {
		AnimeSn  int    `json:"animeSn"`
		Deviceid string `json:"deviceid"`
	}
	resp1, err := core.GET(c, "https://ani.gamer.com.tw/ajax/getdeviceid.php", headers...)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	if len(b1) > 0 && b1[0] == '<' {
		return core.Result{Status: core.StatusNo}
	}
	var res1 apiResponse
	if err := json.Unmarshal(b1, &res1); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp2, err := core.GET(c, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=37783&device="+res1.Deviceid, headers...)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res2 apiResponse
	if err := json.Unmarshal(b2, &res2); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if res2.AnimeSn != 0 {
		resp3, err := core.GET(c, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=38832&device="+res1.Deviceid, headers...)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()
		b3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}

		var res3 apiResponse
		if err := json.Unmarshal(b3, &res3); err != nil {
			return core.Result{Status: core.StatusErr, Err: err}
		}

		if res3.AnimeSn != 0 {
			return core.Result{Status: core.StatusOK, Region: "tw"}
		}

		resp4, err := core.GET(c, "https://ani.gamer.com.tw/cdn-cgi/trace", headers...)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp4.Body.Close()
		b4, err := io.ReadAll(resp4.Body)
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}

		bodyString := string(b4)
		index := strings.Index(bodyString, "loc=")
		if index == -1 {
			return core.Result{Status: core.StatusUnexpected}
		}
		bodyString = bodyString[index+4:]
		index = strings.Index(bodyString, "\n")
		if index == -1 {
			return core.Result{Status: core.StatusUnexpected}
		}
		loc := bodyString[:index]
		if len(loc) == 2 {
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}

