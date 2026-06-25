package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func DirecTVGO(c http.Client) core.Result {
	return core.CheckGETStatus(c, "https://www.directvgo.com/registrarse", core.ResultMap{
		http.StatusForbidden:        {Status: core.StatusNo},
		http.StatusOK:               {Status: core.StatusNo},
		http.StatusMovedPermanently: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}

