package mediaunlocktest

import (
	"net/http"
)

func KayoSports(c http.Client) Result {
	return CheckGETStatus(c, "https://kayosports.com.au/", ResultMap{
		http.StatusForbidden: {Status: StatusNo},
		http.StatusOK:        {Status: StatusOK},
	}, Result{Status: StatusUnexpected},
		H{"Accept", "*/*"},
		H{"Accept-Language", "en-US,en;q=0.9"},
		H{"Origin", "https://kayosports.com.au"},
		H{"Referer", "https://kayosports.com.au/"},
	)
}
