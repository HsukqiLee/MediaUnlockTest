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

func extractWowowDramaURL(body string) []string {
	re := regexp.MustCompile(`https://www.wowow.co.jp/drama/original/\w+/`)
	matches := re.FindAllString(body, -1)

	if len(matches) > 0 {
		return matches
	}
	return []string{}
}

func extractWowowProgramID(body string) string {
	re := regexp.MustCompile(`PRG_CD\s*:\s*"(\d+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func extractWowowMetaID(body string) string {
	re := regexp.MustCompile(`https://wod.wowow.co.jp/watch/(\d+)`)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Wowow(c http.Client) Result {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	resp1, err := GET(c, "https://www.wowow.co.jp/drama/original/json/lineup.json?_="+strconv.FormatInt(timestamp, 10),
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
	//var res1 []struct {
	//    DramaLink string `json:"link"`
	//}
	//if err := json.Unmarshal(body1, &res1); err != nil {
	//	return Result{Status: StatusFailed, Err: err}
	//}

	dramaLinks := extractWowowDramaURL(string(body1))

	if len(dramaLinks) == 0 {
		return Result{Status: StatusFailed}
	}

	var wodUrl string
	for _, dramaLink := range dramaLinks {
		resp2, err := GET(c, dramaLink)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			continue
		}
		programID := extractWowowProgramID(string(body2))
		if programID == "" {
			continue
		}
		resp3, err := PostJson(c, "https://www.wowow.co.jp/API/new_prg/programdetail.php",
			`{"prg_cd": "`+programID+`", "mode": "19"}`,
		)
		if err != nil {
			continue
		}
		var res3 struct {
			ArchiveURL string `json:"archive_url"`
		}
		body3, err := io.ReadAll(resp3.Body)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(body3, &res3); err != nil {
			return Result{Status: StatusErr, Err: err}
		}
		wodUrl = res3.ArchiveURL
		if wodUrl != "" {
			break
		}
	}
	if wodUrl == "" {
		return Result{Status: StatusFailed}
	}

	resp4, err := GET(c, wodUrl)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp4.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	metaID := extractWowowMetaID(string(body4))
	vUid := md5Sum(strconv.FormatInt(timestamp, 10))

	resp5, err := PostJson(c, "https://mapi.wowow.co.jp/api/v1/playback/auth",
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
	defer resp5.Body.Close()

	body5, err := io.ReadAll(resp4.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString5 := string(body5)
	if strings.Contains(bodyString5, "VPN") || resp5.StatusCode == 403 {
		return Result{Status: StatusNo}
	}

	if strings.Contains(bodyString5, "Unauthorized") || strings.Contains(bodyString5, "playback_session_id") || resp5.StatusCode == 401 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
