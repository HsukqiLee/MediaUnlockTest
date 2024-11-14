package mediaunlocktest

import (
	"encoding/json"
	"io"
	"net/http"
)

func NPOStartPlus(c http.Client) Result {
	resp, err := GET(c, "https://npo.nl/start/api/domain/player-token?productId=LI_NL1_4188102",
		H{"connection", "keep-alive"},
		H{"referer", "https://npo.nl/start/live?channel=NPO1"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	var res struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return Result{Status: StatusErr, Err: err}
	}

	resp2, err := PostJson(c, "https://prod.npoplayer.nl/stream-link", `{"profileName":"dash","drmType":"playready","referrerUrl":"https://npo.nl/start/live?channel=NPO1"}`,
		H{"authorization", res.Token},
		H{"referer", "https://npo.nl/"},
		H{"origin", "https://npo.nl"},
	)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}

	if resp2.StatusCode == 451 || resp2.StatusCode == 401 {
		return Result{Status: StatusNo}
	}

	if resp2.StatusCode == 200 {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}
