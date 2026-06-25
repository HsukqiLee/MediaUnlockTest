package mediaunlocktest

import (
	"net/http"
)

func Kancolle(c http.Client) Result {
	return CheckDalvikStatus(c, "https://w00g.kancolle-server.com/kcscontents/news/", ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusForbidden: {Status: StatusNo},
		http.StatusFound:     {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
