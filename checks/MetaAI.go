package mediaunlocktest

import (
	"net/http"
)

func MetaAI(c http.Client) Result {
	return CheckGETStatus(c, "https://www.meta.ai/ajax", ResultMap{
		200: {Status: StatusNo},
		400: {Status: StatusOK},
	}, Result{Status: StatusUnexpected})
}
