package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func AISPlay(c core.HttpClient) core.Result {
	userId := "09e8b25510"
	fakeApiKey := core.MD5Sum(core.GenUUID())
	fakeUdid := core.MD5Sum(core.GenUUID())
	timestamp := time.Now().Unix()

	resp1, err := core.PostJson(c, "https://web-tls.ais-vidnt.com/device/login/?d=gstweb&gst=1&user="+userId+"&pass=e49e9f9e7f", `------WebKitFormBoundaryBj2RhUIW7BtRvfK0--\r\n`,
		core.H{"accept-language", "th"},
		core.H{"api-version", "2.8.2"},
		core.H{"api_key", fakeApiKey},
		core.H{"content-type", "multipart/form-data; boundary=----WebKitFormBoundaryBj2RhUIW7BtRvfK0"},
		core.H{"device-info", "com.vimmi.ais.portal, Windows + Chrome, AppVersion: 4.9.97, 10, language: tha"},
		core.H{"origin", "https://aisplay.ais.co.th"},
		core.H{"privateid", userId},
		core.H{"referer", "https://aisplay.ais.co.th/"},
		core.H{"time", strconv.FormatInt(timestamp, 10)},
		core.H{"udid", fakeUdid},
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
		Info struct {
			Sid string `json:"sid"`
			Dat string `json:"dat"`
		} `json:"info"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}
	sId := res1.Info.Sid
	datAuth := res1.Info.Dat

	timestamp = time.Now().Unix()

	resp2, err := core.GET(c, "https://web-sila.ais-vidnt.com/playtemplate/?d=gstweb",
		core.H{"accept-language", "en-US,en;q=0.9"},
		core.H{"api-version", "2.8.2"},
		core.H{"api_key", fakeApiKey},
		core.H{"dat", datAuth},
		core.H{"device-info", "com.vimmi.ais.portal, Windows + Chrome, AppVersion: 0.0.0, 10, Language: unknown"},
		core.H{"origin", "https://web-player.ais-vidnt.com"},
		core.H{"privateid", userId},
		core.H{"referer", "https://web-player.ais-vidnt.com/"},
		core.H{"sid", sId},
		core.H{"time", strconv.FormatInt(timestamp, 10)},
		core.H{"udid", fakeUdid},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res2 struct {
		Info struct {
			Live string `json:"live"`
		} `json:"Info"`
	}
	if err := json.Unmarshal(body2, &res2); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	mediaId := "B0006"
	realLiveUrl := strings.ReplaceAll(res2.Info.Live, "{MID}", mediaId)
	realLiveUrl = strings.ReplaceAll(realLiveUrl, "metadata.xml", "metadata.json")

	resp3, err := core.GET(c, realLiveUrl+"-https&tuid="+userId+"&tdid="+fakeUdid+"&chunkHttps=true&origin=anevia",
		core.H{"Accept-Language", "en-US,en;q=0.9"},
		core.H{"Origin", "https://web-player.ais-vidnt.com"},
		core.H{"Referer", "https://web-player.ais-vidnt.com/"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp3.Body.Close()

	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res3 struct {
		PlaybackUrls []struct {
			PlayUrl string `json:"url"`
		} `json:"playbackUrls"`
	}
	if err := json.Unmarshal(body3, &res3); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	resp4, err := core.GET(c, res3.PlaybackUrls[0].PlayUrl,
		core.H{"Accept-Language", "en-US,en;q=0.9"},
		core.H{"Origin", "https://web-player.ais-vidnt.com"},
		core.H{"Referer", "https://web-player.ais-vidnt.com/"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	return core.ResultFromMapping(resp4.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
