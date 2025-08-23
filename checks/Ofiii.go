package mediaunlocktest

import (
	"net/http"
)

func Ofiii(c http.Client) Result {
	resp, err := GET(c, "https://cdi.ofiii.com/ofiii_cdi/video/urls?device_type=pc&device_id=b4e377ac-8870-43a4-957a-07f95549a03d&media_type=comic&asset_id=vod68157-020020M001&project_num=OFWEB00&puid=dcafe020-e335-49fb-b9c7-52bd9a15c305")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return ResultFromMapping(resp.StatusCode, ResultMap{
		http.StatusOK:         {Status: StatusOK},
		http.StatusBadRequest: {Status: StatusNo},
	}, Result{Status: StatusUnexpected})
}
