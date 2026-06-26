package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"strings"
)

func EurosportRO(c core.HttpClient) core.Result {
	fakeUuid := core.MD5Sum(core.GenUUID())

	resp1, err := core.GET(c, "https://eu3-prod-direct.eurosport.ro/token?realm=eurosport",
		core.H{"accept", "*/*"},
		core.H{"accept-language", "en-US,en;q=0.9"},
		core.H{"origin", "https://www.eurosport.ro"},
		core.H{"referer", "https://www.eurosport.ro/"},
		core.H{"x-device-info", "escom/0.295.1 (unknown/unknown; Windows/10; " + fakeUuid + ")"},
		core.H{"x-disco-client", "WEB:UNKNOWN:escom:0.295.1"},
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
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	token := res1.Data.Attributes.Token
	sourceSystemId := "eurosport-vid2133403"

	resp2, err := core.GET(c, "https://eu3-prod-direct.eurosport.ro/playback/v2/videoPlaybackInfo/sourceSystemId/"+sourceSystemId+"?usePreAuth=true",
		core.H{"Authorization", "Bearer " + token},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	bodyString := string(body2)

	if strings.Contains(bodyString, "access.denied.geoblocked") {
		return core.Result{Status: core.StatusNo}
	}

	if strings.Contains(bodyString, "eurosport-vod") {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
