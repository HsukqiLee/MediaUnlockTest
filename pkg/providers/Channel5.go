package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Channel5(c http.Client) core.Result {
	timestamp := time.Now().Unix()
	resp, err := core.GET(c, "https://cassie.channel5.com/api/v2/live_media/my5desktopng/C5.json?timestamp="+strconv.FormatInt(timestamp, 10)+"&auth=0_rZDiY0hp_TNcDyk2uD-Kl40HqDbXs7hOawxyqPnbI")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res struct {
		Code string `json:"code"`
	}

	if len(b) > 0 && b[0] == '<' {
		return core.Result{Status: core.StatusBanned}
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusFailed, Err: err}
	}
	switch res.Code {
	case "3000":
		return core.Result{Status: core.StatusNo}
	case "3001":
		return core.Result{Status: core.StatusNo, Info: "Proxy Detected"}
	case "4003":
		return core.Result{Status: core.StatusOK}
	default:
		return core.Result{Status: core.StatusUnexpected}
	}
}

