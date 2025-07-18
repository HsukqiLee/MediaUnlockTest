package mediaunlocktest

import (
	"net/http"
)

func DocPlay(c http.Client) Result {
	return CheckGETStatus(c, "https://www.docplay.com/subscribe", ResultMap{
		http.StatusSeeOther:          {Status: StatusOK},
		http.StatusTemporaryRedirect: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
