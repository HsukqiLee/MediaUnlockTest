package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
	"time"
	"strconv"
	"encoding/json"
)



func AISPlay(c http.Client) Result {
    userId := "09e8b25510"
	fakeApiKey := md5Sum(genUUID())
	fakeUdid := md5Sum(genUUID())
	timestamp := time.Now().Unix() 

	resp1, err := PostJson(c, "https://web-tls.ais-vidnt.com/device/login/?d=gstweb&gst=1&user=" + userId + "&pass=e49e9f9e7f", `------WebKitFormBoundaryBj2RhUIW7BtRvfK0--\r\n`,
	    H{"accept-language", "th"},
	    H{"api-version", "2.8.2"},
	    H{"api_key", fakeApiKey},
	    H{"content-type", "multipart/form-data; boundary=----WebKitFormBoundaryBj2RhUIW7BtRvfK0"},
	    H{"device-info", "com.vimmi.ais.portal, Windows + Chrome, AppVersion: 4.9.97, 10, language: tha"},
	    H{"origin", "https://aisplay.ais.co.th"},
	    H{"privateid", userId},
	    H{"referer", "https://aisplay.ais.co.th/"},
	    H{"time", strconv.FormatInt(timestamp, 10)},
	    H{"udid", fakeUdid},
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
        Info struct {
            Sid string `json:"sid"`
            Dat string `json:"dat"`
        } `json:"info"`
    }
    if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	sId := res1.Info.Sid
	datAuth := res1.Info.Dat
	
	timestamp = time.Now().Unix() 

	resp2, err := GET(c, "https://web-sila.ais-vidnt.com/playtemplate/?d=gstweb",
	    H{"accept-language", "en-US,en;q=0.9"},
	    H{"api-version", "2.8.2"},
	    H{"api_key", fakeApiKey},
	    H{"dat", datAuth},
	    H{"device-info", "com.vimmi.ais.portal, Windows + Chrome, AppVersion: 0.0.0, 10, Language: unknown"},
	    H{"origin", "https://web-player.ais-vidnt.com"},
	    H{"privateid", userId},
	    H{"referer", "https://web-player.ais-vidnt.com/"},
	    H{"sid", sId},
	    H{"time", strconv.FormatInt(timestamp, 10)},
	    H{"udid", fakeUdid},
    )
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

    body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    
    var res2 struct {
        Info struct {
            Live string `json:"live"`
        } `json:"Info"`
    }
    if err := json.Unmarshal(body2, &res2); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	mediaId := "B0006"
	realLiveUrl := strings.ReplaceAll(res2.Info.Live, "{MID}", mediaId)
	realLiveUrl = strings.ReplaceAll(realLiveUrl, "metadata.xml", "metadata.json")


	resp3, err := GET(c, realLiveUrl + "-https&tuid=" + userId + "&tdid=" + fakeUdid + "&chunkHttps=true&origin=anevia",
	    H{"Accept-Language", "en-US,en;q=0.9"},
	    H{"Origin", "https://web-player.ais-vidnt.com"},
	    H{"Referer", "https://web-player.ais-vidnt.com/"},
    )
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()
	
    body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    
    var res3 struct {
        PlaybackUrls []struct {
            PlayUrl string `json:"url"`
        } `json:"playbackUrls"`
    }
    if err := json.Unmarshal(body3, &res3); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	resp4, err:= GET(c, res3.PlaybackUrls[0].PlayUrl,
	    H{"Accept-Language", "en-US,en;q=0.9"},
	    H{"Origin", "https://web-player.ais-vidnt.com"},
	    H{"Referer", "https://web-player.ais-vidnt.com/"},
    )
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()
    
    if resp4.StatusCode == 403 {
        return Result{Status: StatusNo}
    }

    if resp4.StatusCode == 200 {
        return Result{Status: StatusOK}
    }
    
    return Result{Status: StatusUnexpected}
}
 