package mediaunlocktest

import (
	"net/http"
	"strings"
)

func Mora(c http.Client) Result {
	resp, err := GET(c, "https://mora.jp/buy?__requestToken=1713764407153&returnUrl=https%3A%2F%2Fmora.jp%2Fpackage%2F43000087%2FTFDS01006B00Z%2F%3Ffmid%3DTOPRNKS%26trackMaterialNo%3D31168909&fromMoraUx=false&deleteMaterial=")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 || resp.StatusCode == 500 {
		return Result{Status: StatusNo}
	}

	if resp.StatusCode == 302 {
		if strings.Contains(resp.Header.Get("Location"), "error") {
			return Result{Status: StatusNo}
		}
		if strings.Contains(resp.Header.Get("Location"), "signin") {
			return Result{Status: StatusOK}
		}
	}

	return Result{Status: StatusUnexpected}
}
