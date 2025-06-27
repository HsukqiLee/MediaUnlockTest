package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func extractTVerPolicyKey(body string) string {
	re := regexp.MustCompile(`policyKey:"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractTVerDeliveryConfigID(body string) string {
	re := regexp.MustCompile(`deliveryConfigId:"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func isValidTVerEpisodeID(id string) bool {
	re := regexp.MustCompile(`^[a-z0-9]{10}$`)
	return re.MatchString(id)
}

func TVer(c http.Client) Result {
	useDeprecated := false
	if useDeprecated {
		return tver_deprecated(c)
	}
	resp, err := GET(c, "https://playback.api.streaks.jp/v1/projects/tver-simul-ntv/medias/ref:simul-ntv",
		H{"x-streaks-api-key", "ntv"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return Result{Status: StatusOK}
	case 403:
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusUnexpected}
}

func tver_deprecated(c http.Client) Result {
	resp1, err := PostForm(c, "https://platform-api.tver.jp/v2/api/platform_users/browser/create", "device_type=pc",
		H{"origin", "https://s.tver.jp"},
		H{"referer", "https://s.tver.jp/"},
		H{"accept-language", "en-US,en;q=0.9"},
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
		Result struct {
			PlatformUid   string `json:"platform_uid"`
			PlatformToken string `json:"platform_token"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	resp2, err := GET(c, "https://platform-api.tver.jp/service/api/v1/callHome?platform_uid="+res1.Result.PlatformUid+"&platform_token="+res1.Result.PlatformToken+"&require_data=mylist%2Cresume%2Clater",
		H{"origin", "https://tver.jp"},
		H{"referer", "https://tver.jp/"},
		H{"accept-language", "en-US,en;q=0.9"},
		H{"x-tver-platform-type", "web"},
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
		Result struct {
			Components []struct {
				ComponentID string `json:"componentID"`
				Contents    []struct {
					Content struct {
						EpisodeID string `json:"id"`
					} `json:"content"`
				} `json:"contents"`
			} `json:"components"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body2, &res2); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	EpisodeID := ""
	isEpisodeIDValid := false
	for _, component := range res2.Result.Components {
		if component.ComponentID == "variety.catchup.recomend" && len(component.Contents) > 0 {
			for _, content := range component.Contents {
				EpisodeID = content.Content.EpisodeID
				if isValidTVerEpisodeID(EpisodeID) {
					isEpisodeIDValid = true
				}
			}
		}
		if isEpisodeIDValid {
			break
		}
	}
	resp3, err := GET(c, "https://statics.tver.jp/content/episode/"+EpisodeID+".json",
		H{"origin", "https://tver.jp"},
		H{"referer", "https://tver.jp/"},
		H{"accept-language", "en-US,en;q=0.9"},
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
		Video struct {
			AccountID  string `json:"accountID"`
			PlayerID   string `json:"playerID"`
			VideoID    string `json:"videoID"`
			VideoRefID string `json:"videoRefID"`
		} `json:"video"`
	}
	if err := json.Unmarshal(body3, &res3); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}

	AccountID := res3.Video.AccountID
	PlayerID := res3.Video.PlayerID
	VideoID := res3.Video.VideoID
	VideoRefID := res3.Video.VideoRefID

	resp4, err := GET(c, "https://players.brightcove.net/"+AccountID+"/"+PlayerID+"_default/index.min.js",
		H{"Referer", "https://tver.jp/"},
		H{"accept-language", "en-US,en;q=0.9"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp4.Body.Close()

	body4, err := io.ReadAll(resp4.Body)
	if err != nil || len(body4) == 0 {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	PolicyKey := extractTVerPolicyKey(string(body4))
	DeliveryConfigID := extractTVerDeliveryConfigID(string(body4))

	var resp5 *http.Response
	if true { //VideoRefID == "" {
		resp5, err = GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/"+AccountID+"/videos/"+VideoID+"?config_id="+DeliveryConfigID,
			H{"accept", "application/json;pk=" + PolicyKey},
			H{"origin", "https://tver.jp"},
			H{"referer", "https://tver.jp/"},
			H{"accept-language", "en-US,en;q=0.9"},
		)
	} else {
		resp5, err = GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/"+AccountID+"/videos/ref%3A"+VideoRefID,
			H{"accept", "application/json;pk=" + PolicyKey},
			H{"origin", "https://tver.jp"},
			H{"referer", "https://tver.jp/"},
			H{"accept-language", "en-US,en;q=0.9"},
		)
	}
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp5.Body.Close()

	body5, err := io.ReadAll(resp5.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res4a []struct {
		ErrorSubcode string `json:"error_subcode"`
	}
	var res4b struct {
		AccountID string `json:"account_id"`
	}
	if err := json.Unmarshal(body5, &res4a); err != nil {
		if err := json.Unmarshal(body5, &res4b); err != nil {
			return Result{Status: StatusUnexpected}
		}
		if res4b.AccountID != "" {
			return Result{Status: StatusOK}
		}
	}

	if res4a[0].ErrorSubcode == "CLIENT_GEO" {
		return Result{Status: StatusNo}
	}
	return Result{Status: StatusUnexpected}
}
