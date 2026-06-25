package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func BahamutAnime(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)

	headers := getRealisticHeaders("html")
	headers = append(headers, H{"x-custom-headers", "true"})

	type apiResponse struct {
		AnimeSn  int    `json:"animeSn"`
		Deviceid string `json:"deviceid"`
	}
	resp1, err := GET(c, "https://ani.gamer.com.tw/ajax/getdeviceid.php", headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res1 apiResponse
	if err := json.Unmarshal(b1, &res1); err != nil {
		if err.Error() == "invalid character '<' looking for beginning of value" {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusErr, Err: err}
	}

	resp2, err := GET(c, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=37783&device="+res1.Deviceid, headers...)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res2 apiResponse
	if err := json.Unmarshal(b2, &res2); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	if res2.AnimeSn != 0 {
		resp3, err := GET(c, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=38832&device="+res1.Deviceid, headers...)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()
		b3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		var res3 apiResponse
		if err := json.Unmarshal(b3, &res3); err != nil {
			return Result{Status: StatusErr, Err: err}
		}

		if res3.AnimeSn != 0 {
			return Result{Status: StatusOK, Region: "tw"}
		}

		resp4, err := GET(c, "https://ani.gamer.com.tw/cdn-cgi/trace", headers...)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		defer resp4.Body.Close()
		b4, err := io.ReadAll(resp4.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}

		bodyString := string(b4)
		index := strings.Index(bodyString, "loc=")
		if index == -1 {
			return Result{Status: StatusUnexpected}
		}
		bodyString = bodyString[index+4:]
		index = strings.Index(bodyString, "\n")
		if index == -1 {
			return Result{Status: StatusUnexpected}
		}
		loc := bodyString[:index]
		if len(loc) == 2 {
			return Result{Status: StatusOK, Region: strings.ToLower(loc)}
		}
	}
	return Result{Status: StatusUnexpected}
}
