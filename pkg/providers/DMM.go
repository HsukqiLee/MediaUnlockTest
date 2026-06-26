package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func DMM(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://bitcoin.dmm.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "This page is not available in your area") {
		return core.Result{Status: core.StatusNo}
	}
	if strings.Contains(s, "暗号資産") {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo, Info: "Unsupported"}
}

func DMMTV(c core.HttpClient) core.Result {
	resp, err := core.PostJson(c, "https://api.beacon.dmm.com/v1/streaming/start", `{"player_name":"dmmtv_browser","player_version":"0.0.0","content_type_detail":"VOD_SVOD","content_id":"11uvjcm4fw2wdu7drtd1epnvz","purchase_product_id":null}`)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "FOREIGN") {
		return core.Result{Status: core.StatusNo}
	}
	if strings.Contains(s, "UNAUTHORIZED") {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo, Info: "Unsupported"}
}
