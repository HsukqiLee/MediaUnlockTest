package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func MGMPlus(c http.Client) Result {
	resp, err := PostJson(c, "https://api.epix.com/v2/sessions",
		`{"device":{"guid":"7a0baaaf-384c-45cd-a21d-310ca5d3002a","format":"console","os":"web","display_width":1865,"display_height":942,"app_version":"1.0.2","model":"browser","manufacturer":"google"},"apikey":"53e208a9bbaee479903f43b39d7301f7"}`,
		H{"connection", "keep-alive"},
		H{"traceparent", "00-000000000000000015b7efdb572b7bf2-4aefaea90903bd1f-01"},
		H{"x-datadog-sampling-priority", "1"},
		H{"x-datadog-trace-id", "1564983120873880562"},
		H{"x-datadog-parent-id", "5399726519264460063"},
		H{"origin", "https://www.mgmplus.com"},
		H{"referer", "https://www.mgmplus.com/"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	bodyString := string(body)
	if strings.Contains(bodyString, "error code") {
		return Result{Status: StatusNo}
	}
	if strings.Contains(bodyString, "blocked") {
		return Result{Status: StatusBanned}
	}
	var res struct {
		DeviceSession struct {
			SessionToken string `json:"session_token"`
		} `json:"device_session"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	
	resp2, err := PostJson(c, "https://api.epix.com/graphql", `{"operationName":"PlayFlow","variables":{"id":"c2VyaWVzOzEwMTc=","supportedActions":["open_url","show_notice","start_billing","play_content","log_in","noop","confirm_provider","unlinked_provider"],"streamTypes":[{"encryptionScheme":"CBCS","packagingSystem":"DASH"},{"encryptionScheme":"CENC","packagingSystem":"DASH"},{"encryptionScheme":"NONE","packagingSystem":"HLS"},{"encryptionScheme":"SAMPLE_AES","packagingSystem":"HLS"}]},"query":"fragment ShowNotice on ShowNotice {\n  type\n  actions {\n    continuationContext\n    text\n    __typename\n  }\n  description\n  title\n  __typename\n}\n\nfragment OpenUrl on OpenUrl {\n  type\n  url\n  __typename\n}\n\nfragment Content on Content {\n  title\n  __typename\n}\n\nfragment Movie on Movie {\n  id\n  shortName\n  __typename\n}\n\nfragment Episode on Episode {\n  id\n  series {\n    shortName\n    __typename\n  }\n  seasonNumber\n  number\n  __typename\n}\n\nfragment Preroll on Preroll {\n  id\n  __typename\n}\n\nfragment ContentUnion on ContentUnion {\n  ...Content\n  ...Movie\n  ...Episode\n  ...Preroll\n  __typename\n}\n\nfragment PlayContent on PlayContent {\n  type\n  continuationContext\n  heartbeatToken\n  currentItem {\n    content {\n      ...ContentUnion\n      __typename\n    }\n    __typename\n  }\n  nextItem {\n    content {\n      ...ContentUnion\n      __typename\n    }\n    showNotice {\n      ...ShowNotice\n      __typename\n    }\n    showNoticeAt\n    __typename\n  }\n  amazonPlaybackData {\n    pid\n    playbackToken\n    materialType\n    __typename\n  }\n  playheadPosition\n  vizbeeStreamInfo {\n    customStreamInfo\n    __typename\n  }\n  closedCaptions {\n    ttml {\n      location\n      __typename\n    }\n    vtt {\n      location\n      __typename\n    }\n    xml {\n      location\n      __typename\n    }\n    __typename\n  }\n  hints {\n    duration\n    seekAllowed\n    trackingEnabled\n    trackingId\n    __typename\n  }\n  streams(types: $streamTypes) {\n    playlistUrl\n    closedCaptionsEmbedded\n    packagingSystem\n    encryptionScheme\n    videoQuality {\n      height\n      width\n      __typename\n    }\n    widevine {\n      authenticationToken\n      licenseServerUrl\n      __typename\n    }\n    playready {\n      authenticationToken\n      licenseServerUrl\n      __typename\n    }\n    fairplay {\n      authenticationToken\n      certificateUrl\n      licenseServerUrl\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment StartBilling on StartBilling {\n  type\n  __typename\n}\n\nfragment LogIn on LogIn {\n  type\n  __typename\n}\n\nfragment Noop on Noop {\n  type\n  __typename\n}\n\nfragment PreviewContent on PreviewContent {\n  type\n  title\n  description\n  stream {\n    sources {\n      hls {\n        location\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment ConfirmProvider on ConfirmProvider {\n  type\n  __typename\n}\n\nfragment UnlinkedProvider on UnlinkedProvider {\n  type\n  __typename\n}\n\nquery PlayFlow($id: String!, $supportedActions: [PlayFlowActionEnum!]!, $context: String, $behavior: BehaviorEnum = DEFAULT, $streamTypes: [StreamDefinition!]) {\n  playFlow(\n    id: $id\n    supportedActions: $supportedActions\n    context: $context\n    behavior: $behavior\n  ) {\n    ...ShowNotice\n    ...OpenUrl\n    ...PlayContent\n    ...StartBilling\n    ...LogIn\n    ...Noop\n    ...PreviewContent\n    ...ConfirmProvider\n    ...UnlinkedProvider\n    __typename\n  }\n}"}`,
	    H{"x-session-token", res.DeviceSession.SessionToken},
	    H{"connection", "keep-alive"},
		H{"traceparent", "00-000000000000000015b7efdb572b7bf2-4aefaea90903bd1f-01"},
		H{"x-datadog-sampling-priority", "1"},
		H{"x-datadog-trace-id", "1564983120873880562"},
		H{"x-datadog-parent-id", "5399726519264460063"},
		H{"origin", "https://www.mgmplus.com"},
		H{"referer", "https://www.mgmplus.com/"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    bodyString2 := string(body2)
    
    if strings.Contains(bodyString2, "Restricted") {
        return Result{Status: StatusNo}
    }
    if strings.Contains(bodyString2, "StartBilling") {
        return Result{Status: StatusOK}
    }
	
	return Result{Status: StatusUnexpected}
}
