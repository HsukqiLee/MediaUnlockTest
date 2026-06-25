package providers

import (
	"MediaUnlockTest/pkg/core"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func LiTV(c http.Client) core.Result {
	deviceID, err := getLiTVDeviceID(c)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	puid, err := getLiTVPUID(c)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	assetId := "vod70810-000001M001_1500K"
	mediaType := "vod"
	t := time.Now()
	timestamp := t.UnixMilli()
	nonce := genLiTVNonce(t)
	signature := genLiTVSignature(assetId, mediaType, nonce, timestamp)

	payload := map[string]interface{}{
		"AssetId":   assetId,
		"MediaType": mediaType,
		"puid":      puid,
		"timestamp": timestamp,
		"nonce":     nonce,
		"signature": signature,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.PostJson(c, "https://www.litv.tv/api/get-urls-no-auth",
		string(jsonBytes),
		core.H{"Cookie", fmt.Sprintf("device-id=%s; PUID=%s", deviceID, puid)},
		core.H{"Origin", "https://www.litv.tv"},
		core.H{"Referer", "https://www.litv.tv/drama/watch/VOD00328856"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "OutsideRegionError") {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}

func getLiTVDeviceID(c http.Client) (string, error) {
	resp, err := core.PostJson(c, "https://www.litv.tv/api/generate-device-id", "",
		core.H{"Origin", "https://www.litv.tv"},
		core.H{"Referer", "https://www.litv.tv/"},
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "device-id" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("device-id cookie not found")
}

func getLiTVPUID(c http.Client) (string, error) {
	payload := `{"jsonrpc":"2.0","id":100,"method":"PustiService.PUID","params":{"version":"2.0","device_id":"","device_category":"LTWEB00","puid":"","aaid":"","idfa":""}}`
	resp, err := core.PostJson(c, "https://pusti.svc.litv.tv/puid", payload,
		core.H{"Origin", "https://www.litv.tv"},
		core.H{"Referer", "https://www.litv.tv/"},
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Result struct {
			Puid string `json:"puid"`
		} `json:"result"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	if res.Result.Puid == "" {
		return "", fmt.Errorf("puid not found in response")
	}
	return res.Result.Puid, nil
}

func genLiTVNonce(t time.Time) string {
	return core.GenBase36(13) + core.GenBase36(13) + strconv.FormatInt(t.UnixMilli(), 36)
}

func genLiTVSignature(assetId, mediaType, nonce string, timestamp int64) string {
	key := "7f4a9c2e8b6d1f3a5e9c7b4d2f8a6e1c"
	// e + t + r + n + i
	data := assetId + mediaType + strconv.FormatInt(timestamp, 10) + nonce + key
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

