package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Ofiii(c http.Client) core.Result {
	resp, err := core.GET(c, "https://cdi.ofiii.com/ofiii_cdi/video/urls?device_type=pc&device_id=b4e377ac-8870-43a4-957a-07f95549a03d&media_type=comic&asset_id=vod68157-020020M001&project_num=OFWEB00&puid=dcafe020-e335-49fb-b9c7-52bd9a15c305")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:         {Status: core.StatusOK},
		http.StatusBadRequest: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

