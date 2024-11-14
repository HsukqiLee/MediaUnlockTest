package mediaunlocktest

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func urlEncode(input string) string {
	return url.QueryEscape(input)
}

func generateHMACSignature(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func NaverTV(c http.Client) Result {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	signature := generateHMACSignature(
		"nbxvs5nwNG9QKEWK0ADjYA4JZoujF4gHcIwvoCxFTPAeamq5eemvt5IWAYXxrbYM",
		"https://apis.naver.com/now_web2/now_web_api/v1/clips/31030608/play-info"+strconv.FormatInt(timestamp, 10),
	)

	resp, err := GET(c, "https://apis.naver.com/now_web2/now_web_api/v1/clips/31030608/play-info?msgpad="+strconv.FormatInt(timestamp, 10)+"&md="+urlEncode(signature),
		H{"Connection", "keep-alive"},
		H{"Accept", "application/json, text/plain, */*"},
		H{"Origin", "https://tv.naver.com"},
		H{"Referer", "https://tv.naver.com/v/31030608"},
	)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res struct {
		Playable string `json:"playable"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	if res.Playable == "NOT_COUNTRY_AVAILABLE" {
		return Result{Status: StatusNo}
	}
	if resp.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
