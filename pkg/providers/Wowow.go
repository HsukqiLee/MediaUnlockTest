package providers

import (
	"MediaUnlockTest/pkg/core"
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

func Wowow(c http.Client) core.Result {
	useDeprecated := false
	if useDeprecated {
		return wowow_deprecated(c)
	}
	resp, err := core.PostJson(c, "https://mapi.wowow.co.jp/api/v1/playback/auth", `{"meta_id":81174}`)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 403 {
		var res struct {
			Error struct {
				Code int `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return core.Result{Status: core.StatusErr, Err: err}
		}
		switch res.Error.Code {
		case 2055:
			return core.Result{Status: core.StatusNo}
		case 2041, 2003:
			return core.Result{Status: core.StatusOK}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}

func wowow_deprecated(c http.Client) core.Result {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	resp1, err := core.GET(c, "https://www.wowow.co.jp/assets/config/top_recommend_list.json",
		core.H{"Accept", "application/json, text/javascript, */*; q=0.01"},
		core.H{"Referer", "https://www.wowow.co.jp/drama/original/"},
		core.H{"X-Requested-With", "XMLHttpRequest"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res1 struct {
		Movie         []interface{} `json:"movie"`
		DramaOriginal []interface{} `json:"drama_original"`
		Music         []interface{} `json:"music"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	var wodUrl string
	for _, id := range res1.DramaOriginal {
		var programID string
		if str, ok := id.(string); ok {
			programID = str
		} else if num, ok := id.(float64); ok {
			programID = strconv.FormatFloat(num, 'f', 0, 64)
		}
		resp2, err := core.PostJson(c, "https://www.wowow.co.jp/API/new_prg/programdetail.php",
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
			return core.Result{Status: core.StatusErr, Err: err}
		}
		wodUrl = res2.ArchiveURL
		if wodUrl != "" {
			break
		}
	}
	if wodUrl == "" {
		return core.Result{Status: core.StatusFailed}
	}

	resp3, err := core.GET(c, wodUrl)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	metaID := extractWowowMetaID(string(body3))
	vUid := core.MD5Sum(strconv.FormatInt(timestamp, 10))

	resp4, err := core.PostJson(c, "https://mapi.wowow.co.jp/api/v1/playback/auth",
		`{"meta_id":`+metaID+`,"vuid":"`+vUid+`","device_code":1,"app_id":1,"ua":"`+core.UA_Browser+`"}`,
		core.H{"accept", "application/json, text/plain, */*"},
		core.H{"content-type", "application/json;charset=UTF-8"},
		core.H{"origin", "https://wod.wowow.co.jp"},
		core.H{"referer", "https://wod.wowow.co.jp/"},
		core.H{"x-requested-with", "XMLHttpRequest"},
	)

	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp3.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString4 := string(body4)
	if strings.Contains(bodyString4, "VPN") || resp4.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}

	if resp4.StatusCode == 201 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}


