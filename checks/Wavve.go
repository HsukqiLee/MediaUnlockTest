package mediaunlocktest

import (
	"net/http"
)

func Wavve(c http.Client) Result {
	resp, err := GET(c, "https://apis.wavve.com/fz/streaming?device=pc&partner=pooq&apikey=E5F3E0D30947AA5440556471321BB6D9&credential=none&service=wavve&pooqzone=none&region=kor&drm=pr&targetage=all&contentid=MV_C3001_C300000012559&contenttype=movie&hdr=sdr&videocodec=avc&audiocodec=ac3&issurround=n&format=normal&withinsubtitle=n&action=dash&protocol=dash&quality=auto&deviceModelId=Windows%2010&guid=1a8e9c88-6a3b-11ed-8584-eed06ef80652&lastplayid=none&authtype=cookie&isabr=y&ishevc=n")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return Result{Status: StatusBanned}
	}
	
	if resp.StatusCode == 550 || resp.StatusCode == 421 {
		return Result{Status: StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return Result{Status: StatusOK}
	}

	return Result{Status: StatusUnexpected}
}