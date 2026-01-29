package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func AETV(c http.Client) Result {
	// Step 1: Request stream manifest
	step1Url := "https://dai.google.com/ondemand/hls/content/2540935/vid/2400478275868/streams"
	step1Data := `dai-dlid=default_delivery_60101&afid=59946479&adobe_id=17218268117322661452685981206066538688&asnw=171213&caid=305489&imafw__fw_ae=nonauthenticated&imafw__fw_vcid2=2395239744222918&devicename=Desktop&imafw_csid=aetv.desktop.video&imafw__fw_player_height=901&imafw__fw_player_width=474&imafw__fw_us_privacy=1---&imafw__fw_site_page=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&ltd=1&imafw__fw_h_user_agent=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F144.0.0.0%20Safari%2F537.36%20Edg%2F144.0.0.0&imafw_prof=171213%3Aaetn_desktop_ae_googlessai_vod&metr=1031&nw=171213&pvrn=7469863039993294&ssnw=171213&vprn=7469863039993294&flag=%2Bsltp%2Bvicb%2Bslcb%2Bamsl%2Bamcb%2Bssus%2Bemcr%2Bfbad%2Bdtrd%2Bplay&resp=vmap1&cld=1&ctv=0&correlator=3412237903465914&ptt=20&osd=2&sdr=1&sdki=41&sdkv=h.3.740.0&uach=WyJXaW5kb3dzIiwiMTUuMC4wIiwieDg2IiwiIiwiMTQ0LjAuMzcxOS45MiIsbnVsbCwwLG51bGwsIjY0IixbWyJOb3QoQTpCcmFuZCIsIjguMC4wLjAiXSxbIkNocm9taXVtIiwiMTQ0LjAuNzU1OS45NyJdLFsiTWljcm9zb2Z0IEVkZ2UiLCIxNDQuMC4zNzE5LjkyIl1dLDBd&ua=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F144.0.0.0%20Safari%2F537.36%20Edg%2F144.0.0.0&eid=44751890%2C95322027%2C95331589%2C95332046&frm=0&omid_p=Google1%2Fh.3.740.0&sdk_apis=7&sid=6C436384-F48F-4AEA-8FCA-9F3633757A74&ssss=gima&ref=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&url=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&wta=0&us_privacy=1---&gpp_sid=-1&eoidce=1`

	req1, err := http.NewRequest("POST", step1Url, strings.NewReader(step1Data))
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	req1.Header.Set("accept", "*/*")
	req1.Header.Set("accept-language", "zh-HK,zh;q=0.9")
	req1.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req1.Header.Set("origin", "https://play.aetv.com")
	req1.Header.Set("preferanonymous", "1")
	req1.Header.Set("priority", "u=1, i")
	req1.Header.Set("referer", "https://play.aetv.com/")
	req1.Header.Set("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Microsoft Edge";v="144"`)
	req1.Header.Set("sec-ch-ua-mobile", "?0")
	req1.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req1.Header.Set("sec-fetch-dest", "empty")
	req1.Header.Set("sec-fetch-mode", "cors")
	req1.Header.Set("sec-fetch-site", "cross-site")
	req1.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0")

	// Use cdo for retries
	resp1, err := cdo(c, req1)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != 200 && resp1.StatusCode != 201 {
		return Result{Status: StatusUnexpected}
	}

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}

	var res1 struct {
		StreamManifest string `json:"stream_manifest"`
	}
	if err := json.Unmarshal(body1, &res1); err != nil {
		return Result{Status: StatusUnexpected, Err: err}
	}

	if res1.StreamManifest == "" {
		return Result{Status: StatusUnexpected, Info: "No stream manifest found"}
	}

	// Step 2: Check access
	manifestUrl, err := url.Parse(res1.StreamManifest)
	if err != nil {
		return Result{Status: StatusUnexpected, Err: err}
	}
	originPath := manifestUrl.Query().Get("originpath")
	if originPath == "" {
		// Fallback to trying to extract from path if query is empty (though example shows checks on query)
		// For now assume strictly it's in originpath
	}

	// Construct Step 2 URL
	base2 := "https://argus.media.aetnd.com/"
	u2, _ := url.Parse(base2)
	q2 := u2.Query()
	q2.Set("bestCdn", "fastly")
	q2.Set("brand", "aetv")
	q2.Set("dfpOp", originPath)
	q2.Set("client", "tve-web-theo")
	q2.Set("sig", "00697b29892c60831e3de250beaec91bd6ab73901a0ffb723b533130425058484d6c62")
	u2.RawQuery = q2.Encode()

	req2, err := http.NewRequest("GET", u2.String(), nil)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	req2.Header.Set("accept", "*/*")
	req2.Header.Set("accept-language", "zh-HK,zh;q=0.9")
	req2.Header.Set("origin", "https://play.aetv.com")
	req2.Header.Set("priority", "u=1, i")
	req2.Header.Set("referer", "https://play.aetv.com/")
	req2.Header.Set("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Microsoft Edge";v="144"`)
	req2.Header.Set("sec-ch-ua-mobile", "?0")
	req2.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req2.Header.Set("sec-fetch-dest", "empty")
	req2.Header.Set("sec-fetch-mode", "cors")
	req2.Header.Set("sec-fetch-site", "cross-site")
	req2.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0")
	req2.Header.Set("x-video-meta-token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3Njk2NzU2NDAsImV4cCI6MTc2OTY3OTI0MCwicHBsSWQiOiIzMDU0ODkiLCJpc0JlaGluZFdhbGwiOmZhbHNlLCJpc0xvbmdGb3JtIjp0cnVlLCJlbmNvZGVyIjoiYml0bW92aW5fdjEiLCJyZW5kaXRpb25zTGlzdCI6bnVsbCwicmVuZGl0aW9uc1BhdGhQcmVmaXgiOiJiaXRtb3Zpbi9BRVROLUFFVFZfVk1TL0FFTl9PWlJLXzMwNTQ4OV9HTEJfNDkzNjY5XzIzOThfNjBfMjAyNTAxMDRfMDFfQUVUTi1BRVRWX1ZNUyIsInJlZ2lvbnNBdmFpbGFibGUiOlsiVVMiLCJDQSIsIkFTIiwiR1UiLCJNUCIsIlBSIiwiVkkiLCJVTSJdfQ._IBIJ-Yh8X9hGOGglCNQWcW6WZKrUbjTNQ8Rdon8u_A")

	resp2, err := cdo(c, req2)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 403 {
		b, err := io.ReadAll(resp2.Body)
		if err != nil {
			return Result{Status: StatusNetworkErr, Err: err}
		}
		var res2 struct {
			Exception string `json:"exception"`
		}
		if err := json.Unmarshal(b, &res2); err != nil {
			return Result{Status: StatusUnexpected, Err: err}
		}
		if res2.Exception == "GeoLocationBlocked" {
			return Result{Status: StatusNo}
		}
		if res2.Exception == "JWTExpiredSignature" {
			return Result{Status: StatusOK}
		}
	}

	return Result{Status: StatusUnexpected, Info: "Step 2 status: " + resp2.Status}
}
