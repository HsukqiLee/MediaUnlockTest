package mediaunlocktest

import (
	"io"
	"net/http"
	"strings"
)

func SupportApple(loc string) bool {
	var APPLE_SUPPORT_COUNTRY = []string{
		"ALB", "DZA", "AFG", "ARG", "ARE", "ABW", "OMN", "AZE", "EGY", "ETH", "IRL", "EST", "AND", "AGO", "AIA", "ATG",
		"AUT", "ALA", "AUS", "MAC", "BRB", "PNG", "BHS", "PAK", "PRY", "PSE", "BHR", "PAN", "BRA", "BLR", "BMU", "BGR",
		"MNP", "MKD", "BEN", "BEL", "ISL", "BOL", "PRI", "POL", "BIH", "BWA", "BLZ", "BTN", "BFA", "BDI", "BVT", "IOT",
		"GNQ", "DNK", "DEU", "TLS", "TGO", "DOM", "DMA", "RUS", "ECU", "ERI", "FRA", "FRO", "PYF", "GUF", "ATF", "PHL",
		"FIN", "CPV", "FLK", "GMB", "COG", "COD", "COL", "CRI", "GRD", "GRL", "GEO", "GGY", "GLP", "GUM", "GUY", "KAZ",
		"HTI", "KOR", "NLD", "BES", "SXM", "HMD", "MNE", "HND", "KIR", "DJI", "KGZ", "GIN", "GNB", "CAN", "GHA", "GAB",
		"KHM", "CZE", "ZWE", "CMR", "QAT", "CYM", "CCK", "COM", "XKS", "CIV", "KWT", "HRV", "KEN", "COK", "CUW", "LVA",
		"LSO", "LAO", "LBN", "LBR", "LBY", "LTU", "LIE", "REU", "LUX", "RWA", "ROU", "MDG", "MLT", "MDV", "MWI", "MYS",
		"MLI", "MHL", "MTQ", "MYT", "IMN", "MUS", "MRT", "USA", "UMI", "ASM", "VIR", "MNG", "MSR", "BGD", "PER", "FSM",
		"MMR", "MDA", "MAR", "MCO", "MOZ", "MEX", "NAM", "ZAF", "ATA", "SGS", "SSD", "NPL", "NIC", "NER", "NGA", "NIU",
		"NOR", "NFK", "PLW", "PCN", "PRT", "JPN", "SWE", "CHE", "SLV", "WSM", "SRB", "SLE", "SEN", "CYP", "SYC", "SAU",
		"BLM", "CXR", "STP", "SHN", "KNA", "LCA", "MAF", "SMR", "SPM", "VCT", "LKA", "SVK", "SVN", "SJM", "SWZ", "SDN",
		"SUR", "SOM", "SLB", "TJK", "TWN", "THA", "TZA", "TON", "TCA", "TTO", "TUN", "TUV", "TUR", "TKM", "TKL", "WLF",
		"VUT", "GTM", "VEN", "BRN", "UGA", "UKR", "URY", "UZB", "ESP", "ESH", "GRC", "HKG", "SGP", "NCL", "NZL", "HUN",
		"JAM", "ARM", "YEM", "IRQ", "ISR", "ITA", "IND", "IDN", "GBR", "VGB", "JOR", "VNM", "ZMB", "JEY", "TCD", "GIB",
		"CHL", "CAF", "CHN", "NRU", "VAT", "FJI",
	}
	for _, s := range APPLE_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func Apple(c http.Client) Result {
	resp, err := GET(c, "https://gspe1-ssl.ls.apple.com/pep/gcc")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	s := string(b)
	loc := twoToThreeCode(s)

	if SupportApple(loc) {
		return Result{Status: StatusOK, Region: strings.ToLower(s)}
	}
	return Result{Status: StatusNo, Region: strings.ToLower(s)}
}
