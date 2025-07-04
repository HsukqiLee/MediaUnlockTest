package mediaunlocktest

import (
    "net/http"
    "strings"
    "io"
)

func Sky_DE(c http.Client) Result {
    resp, err := GET(c, "https://edge.api.brightcove.com/playback/v1/accounts/1050888051001/videos/6247131490001",
        H{"Accept", "application/json;pk=BCpkADawqM0OXCLe4eIkpyuir8Ssf3kIQAM62a1KMa4-1_vTOWQIxoHHD4-oL-dPmlp-rLoS-WIAcaAMKuZVMR57QY4uLAmP4Ov3V416hHbqr0GNNtzVXamJ6d4-rA3Xi98W-8wtypdEyjGEZNepUCt3D7UdMthbsG-Ean3V4cafT4nZX03st5HlyK1chp51SfA-vKcAOhHZ4_Oa9TTN61tEH6YqML9PWGyKrbuN5myICcGsFzP3R2aOF8c5rPCHT2ZAiG7MoavHx8WMjhfB0QdBr2fphX24CSpUKlcjEnQJnBiA1AdLg9iyReWrAdQylX4Eyhw5OwKiCGJznfgY6BDtbUmeq1I9r9RfmhP5bfxVGjILSEFZgXbMqGOvYdrdare0aW2fTCxeHdHt0vyKOWTC6CS1lrGJF2sFPKn1T1csjVR8s4MODqCBY1PTbHY4A9aZ-2MDJUVJDkOK52hGej6aXE5b9N9_xOT2B9wbXL1B1ZB4JLjeAdBuVtaUOJ44N0aCd8Ns0o02E1APxucQqrjnEociLFNB0Bobe1nkGt3PS74IQcs-eBvWYSpolldMH6TKLu8JqgdnM4WIp3FZtTWJRADgAmvF9tVDUG9pcJoRx_CZ4im-rn-AzN3FeOQrM4rTlU3Q8YhSmyEIoxYYqsFDwbFlhsAcvqQkgaElYtuciCL5i3U8N4W9rIhPhQJzsPafmLdWxBP_FXicyek25GHFdQzCiT8nf1o860Jv2cHQ4xUNcnP-9blIkLy9JmuB2RgUXOHzWsrLGGW6hq9wLUtqwEoxcEAAcNJgmoC0k8HE-Ga-NHXng6EFWnqiOg_mZ_MDd7gmHrrKLkQV"},
        H{"Origin", "https://www.sky.de"},
    )
    if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
    
    if strings.Contains(string(b), "CLIENT_GEO") {
        return Result{Status: StatusNo}
    }
    
    if resp.StatusCode == 200 {
        return Result{Status: StatusOK}
    }
    
    return Result{Status: StatusUnexpected}
}
