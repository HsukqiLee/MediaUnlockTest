package providers

import (
	"MediaUnlockTest/pkg/core"
	"encoding/json"
	"io"
)

func EncoreTVB(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/5324042807001/videos/6005570109001",
		core.H{"Accept", "application/json;pk=BCpkADawqM2Gpjj8SlY2mj4FgJJMfUpxTNtHWXOItY1PvamzxGstJbsgc-zFOHkCVcKeeOhPUd9MNHEGJoVy1By1Hrlh9rOXArC5M5MTcChJGU6maC8qhQ4Y8W-QYtvi8Nq34bUb9IOvoKBLeNF4D9Avskfe9rtMoEjj6ImXu_i4oIhYS0dx7x1AgHvtAaZFFhq3LBGtR-ZcsSqxNzVg-4PRUI9zcytQkk_YJXndNSfhVdmYmnxkgx1XXisGv1FG5GOmEK4jZ_Ih0riX5icFnHrgniADr4bA2G7TYh4OeGBrYLyFN_BDOvq3nFGrXVWrTLhaYyjxOr4rZqJPKK2ybmMsq466Ke1ZtE-wNQ"},
		core.H{"Origin", "https://www.encoretvb.com"},
	)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	var res1 struct {
		ErrorSubcode string `json:"error_subcode"`
		AccountId    string `json:"account_id"`
	}
	var res2 []struct {
		ErrorSubcode string `json:"error_subcode"`
		//ClientGeo    string `json:"client_geo"`
	}
	if err := json.Unmarshal(b, &res1); err != nil {
		if err := json.Unmarshal(b, &res2); err != nil {
			return core.Result{Status: core.StatusFailed, Err: err}
		}
		if res2[0].ErrorSubcode == "CLIENT_GEO" {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusFailed, Err: err}
	}

	if res1.AccountId != "0" {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusUnexpected}
}
