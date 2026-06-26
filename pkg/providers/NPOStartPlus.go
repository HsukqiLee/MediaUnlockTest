package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func NPOStartPlus(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://npo.nl/start/api/domain/player-token?productId=LI_NL1_4188102",
		core.H{"connection", "keep-alive"},
		core.H{"referer", "https://npo.nl/start/live?channel=NPO1"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	var res struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp2, err := core.PostJson(c, "https://prod.npoplayer.nl/stream-link", `{"profileName":"dash","drmType":"playready","referrerUrl":"https://npo.nl/start/live?channel=NPO1"}`,
		core.H{"authorization", res.Token},
		core.H{"referer", "https://npo.nl/"},
		core.H{"origin", "https://npo.nl"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp2.StatusCode == 451 || resp2.StatusCode == 401 {
		return core.Result{Status: core.StatusNo}
	}

	if resp2.StatusCode == 200 {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
