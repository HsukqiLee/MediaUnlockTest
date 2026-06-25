package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Kancolle(c http.Client) core.Result {
	return core.CheckDalvikStatus(c, "https://w00g.kancolle-server.com/kcscontents/news/", core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusFound:     {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

