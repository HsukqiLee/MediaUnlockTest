package mediaunlocktest

import (
	"encoding/json"
	"io"
	"strings"
	"net/http"
	"net/http/cookiejar"
)

func BahamutAnime(c http.Client) Result {
	c.Jar, _ = cookiejar.New(nil)
	resp, err := GET(c, "https://ani.gamer.com.tw/ajax/getdeviceid.php")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		AnimeSn  int
		Deviceid string
	}
	if err := json.Unmarshal(b, &res); err != nil {
	    if err.Error() == "invalid character '<' looking for beginning of value" {
	        return Result{Status: StatusNo}
	    }
		return Result{Status: StatusErr, Err: err}
	}
	resp, err = GET(c, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667&device="+res.Deviceid)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}
	if res.AnimeSn != 0 {
	    resp, err = GET(c, "https://ani.gamer.com.tw/cdn-cgi/trace")
	    if err != nil {
		    return Result{Status: StatusNetworkErr, Err: err}
	    }
	    defer resp.Body.Close()
	    b, err = io.ReadAll(resp.Body)
	    if err != nil {
	    	return Result{Status: StatusNetworkErr, Err: err}
	    }
	    bodyString := string(b)
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
