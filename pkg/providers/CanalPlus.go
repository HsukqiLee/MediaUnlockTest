package mediaunlocktest

import (
	"net/http"
)

func CanalPlus(c http.Client) Result {
	return CheckGETStatus(c, "https://boutique-tunnel.canalplus.com/", ResultMap{
		http.StatusOK:        {Status: StatusOK},
		http.StatusFound:     {Status: StatusNo},
		http.StatusForbidden: {Status: StatusBanned},
	}, Result{Status: StatusUnexpected})
}
