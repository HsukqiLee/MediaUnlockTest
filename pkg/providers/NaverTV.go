package providers

import (
	"MediaUnlockTest/pkg/core"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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

func NaverTV(c core.HttpClient) core.Result {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	signature := generateHMACSignature(
		"nbxvs5nwNG9QKEWK0ADjYA4JZoujF4gHcIwvoCxFTPAeamq5eemvt5IWAYXxrbYM",
		"https://apis.naver.com/now_web2/now_web_api/v1/clips/31030608/play-info"+strconv.FormatInt(timestamp, 10),
	)

	resp, err := core.GET(c, "https://apis.naver.com/now_web2/now_web_api/v1/clips/31030608/play-info?msgpad="+strconv.FormatInt(timestamp, 10)+"&md="+urlEncode(signature),
		core.H{"Connection", "keep-alive"},
		core.H{"Accept", "application/json, text/plain, */*"},
		core.H{"Origin", "https://tv.naver.com"},
		core.H{"Referer", "https://tv.naver.com/v/31030608"},
	)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Playable string `json:"playable"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if res.Playable == "NOT_COUNTRY_AVAILABLE" {
		return core.Result{Status: core.StatusNo}
	}
	if resp.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
