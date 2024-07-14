package mediaunlocktest

import (
    "time"
    "regexp"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
)

func extractWowowDramaURL(body string) string {
    re := regexp.MustCompile(`https://www.wowow.co.jp/drama/original/\w+/`)
    matches := re.FindStringSubmatch(body)
    
    if len(matches) > 0 {
        return matches[0]
    }
    return ""
}

func extractWowowContentURL(body string) string {
    re := regexp.MustCompile(`https://wod.wowow.co.jp/content/\d+`)
    match := re.FindString(body)
    return match
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
    timestamp := time.Now().UnixNano()/int64(time.Millisecond)
    resp1, err := GET(c, "https://www.wowow.co.jp/drama/original/json/lineup.json?_=" + strconv.FormatInt(timestamp, 10),
        H{"Accept", "application/json, text/javascript, */*; q=0.01"},
        H{"Referer", "https://www.wowow.co.jp/drama/original/"},
        H{"X-Requested-With", "XMLHttpRequest"},
    )
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    defer resp1.Body.Close()

    body1, err := ioutil.ReadAll(resp1.Body)
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    //var res1 []struct {
    //    DramaLink string `json:"link"`
    //}
    //if err := json.Unmarshal(body1, &res1); err != nil {
	//	return Result{Status: StatusFailed, Err: err}
	//}

	dramaLink := extractWowowDramaURL(string(body1))
	
	if dramaLink == "" {
	    return Result{Status: StatusFailed}
	}

    resp2, err := GET(c, dramaLink)
    //resp2, err := GET(c, "https://www.wowow.co.jp/drama/original/yukai/")
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    defer resp2.Body.Close()

    body2, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    
    wodUrl := extractWowowContentURL(string(body2))
    if wodUrl == "" {
        return Result{Status: StatusFailed, Err: err}
    }
    
    resp3, err := GET(c, wodUrl)
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    defer resp3.Body.Close()

    body3, err := ioutil.ReadAll(resp3.Body)
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }

    metaID := extractWowowMetaID(string(body3))
    vUid := md5Sum(strconv.FormatInt(timestamp, 10) )

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

    body4, err := ioutil.ReadAll(resp4.Body)
    if err != nil {
        return Result{Status: StatusNetworkErr, Err: err}
    }
    bodyString4 := string(body4)
    
    if strings.Contains(bodyString4, "VPN") {
        return Result{Status: StatusNo}
    }
    
    if strings.Contains(bodyString4, "Unauthorized") || strings.Contains(bodyString4, "playback_session_id") {
        return Result{Status: StatusOK}
    }

    return Result{Status: StatusUnexpected}
}
