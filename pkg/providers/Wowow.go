package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func extractWowowMetaID(body string) string {
	re := regexp.MustCompile(`https://wod.wowow.co.jp/watch/(\d+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Wowow(c http.Client) Result {
	useDeprecated := false
	if useDeprecated {
		return wowow_deprecated(c)
	}
	resp, err := PostJson(c, "https://mapi.wowow.co.jp/api/v1/playback/auth", `{"meta_id":81174}`)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 403 {
		var res struct {
			Error struct {
				Code int `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return Result{Status: StatusErr, Err: err}
		}
		switch res.Error.Code {
		case 2055:
			return Result{Status: StatusNo}
		case 2041, 2003:
			return Result{Status: StatusOK}
		}
	}
	return Result{Status: StatusUnexpected}
}

func wowow_deprecated(c http.Client) Result {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	resp1, err := GET(c, "https://www.wowow.co.jp/assets/config/top_recommend_list.json",
		H{"Accept", "application/json, text/javascript, */*; q=0.01"},
		H{"Referer", "https://www.wowow.co.jp/drama/original/"},
		H{"X-Requested-With", "XMLHttpRequest"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res1 struct {
		Movie         []interface{} `json:"movie"`
		DramaOriginal []interface{} `json:"drama_original"`
		Music         []interface{} `json:"music"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	var wodUrl string
	for _, id := range res1.DramaOriginal {
		var programID string
		if str, ok := id.(string); ok {
			programID = str
		} else if num, ok := id.(float64); ok {
			programID = strconv.FormatFloat(num, 'f', 0, 64)
		}
		resp2, err := PostJson(c, "https://www.wowow.co.jp/API/new_prg/programdetail.php",
			`{"prg_cd": "`+programID+`", "mode": "19"}`,
		)
		if err != nil {
			continue
		}
		var res2 struct {
			ArchiveURL string `json:"archive_url"`
		}
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(body2, &res2); err != nil {
			return Result{Status: StatusErr, Err: err}
		}
		wodUrl = res2.ArchiveURL
		if wodUrl != "" {
			break
		}
	}
	if wodUrl == "" {
		return Result{Status: StatusFailed}
	}

	resp3, err := GET(c, wodUrl)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	metaID := extractWowowMetaID(string(body3))
	vUid := md5Sum(strconv.FormatInt(timestamp, 10))

	resp4, err := PostJson(c, "https://mapi.wowow.co.jp/api/v1/playback/auth",
		`{"meta_id":`+metaID+`,"vuid":"`+vUid+`","device_code":1,"app_id":1,"ua":"`+UA_Browser+`"}`,
		H{"accept", "application/json, text/plain, */*"},
		H{"content-type", "application/json;charset=UTF-8"},
		H{"origin", "https://wod.wowow.co.jp"},
		H{"referer", "https://wod.wowow.co.jp/"},
		H{"x-requested-with", "XMLHttpRequest"},
	)

	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp3.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString4 := string(body4)
	if strings.Contains(bodyString4, "VPN") || resp4.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if resp4.StatusCode == 201 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
