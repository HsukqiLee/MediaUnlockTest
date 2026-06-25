package mediaunlocktest

import (
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

func LiTV(c http.Client) Result {
	deviceID, err := getLiTVDeviceID(c)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	puid, err := getLiTVPUID(c)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
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
		return Result{Status: StatusErr, Err: err}
	}

	resp, err := PostJson(c, "https://www.litv.tv/api/get-urls-no-auth",
		string(jsonBytes),
		H{"Cookie", fmt.Sprintf("device-id=%s; PUID=%s", deviceID, puid)},
		H{"Origin", "https://www.litv.tv"},
		H{"Referer", "https://www.litv.tv/drama/watch/VOD00328856"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "OutsideRegionError") {
			return Result{Status: StatusNo}
		}
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}

func getLiTVDeviceID(c http.Client) (string, error) {
	resp, err := PostJson(c, "https://www.litv.tv/api/generate-device-id", "",
		H{"Origin", "https://www.litv.tv"},
		H{"Referer", "https://www.litv.tv/"},
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
	resp, err := PostJson(c, "https://pusti.svc.litv.tv/puid", payload,
		H{"Origin", "https://www.litv.tv"},
		H{"Referer", "https://www.litv.tv/"},
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
	return genBase36(13) + genBase36(13) + strconv.FormatInt(t.UnixMilli(), 36)
}

func genLiTVSignature(assetId, mediaType, nonce string, timestamp int64) string {
	key := "7f4a9c2e8b6d1f3a5e9c7b4d2f8a6e1c"
	// e + t + r + n + i
	data := assetId + mediaType + strconv.FormatInt(timestamp, 10) + nonce + key
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
