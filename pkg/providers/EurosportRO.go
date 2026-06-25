package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func EurosportRO(c http.Client) Result {
	fakeUuid := md5Sum(genUUID())

	resp1, err := GET(c, "https://eu3-prod-direct.eurosport.ro/token?realm=eurosport",
		H{"accept", "*/*"},
		H{"accept-language", "en-US,en;q=0.9"},
		H{"origin", "https://www.eurosport.ro"},
		H{"referer", "https://www.eurosport.ro/"},
		H{"x-device-info", "escom/0.295.1 (unknown/unknown; Windows/10; " + fakeUuid + ")"},
		H{"x-disco-client", "WEB:UNKNOWN:escom:0.295.1"},
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
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	token := res1.Data.Attributes.Token
	sourceSystemId := "eurosport-vid2133403"

	resp2, err := GET(c, "https://eu3-prod-direct.eurosport.ro/playback/v2/videoPlaybackInfo/sourceSystemId/"+sourceSystemId+"?usePreAuth=true",
		H{"Authorization", "Bearer " + token},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString := string(body2)

	if strings.Contains(bodyString, "access.denied.geoblocked") {
		return Result{Status: StatusNo}
	}

	if strings.Contains(bodyString, "eurosport-vod") {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
