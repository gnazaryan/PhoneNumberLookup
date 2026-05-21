package phone

import "strings"

var isoToDialing = map[string]string{
	"US": "1", "CA": "1",
	"MX": "52", "ES": "34", "PT": "351", "GB": "44",
	"DE": "49", "FR": "33", "IT": "39",
	"BR": "55", "AR": "54",
	"JP": "81", "CN": "86", "IN": "91",
	"AU": "61", "NZ": "64",
}

var dialingToISO = map[string]string{
	"1":   "US",
	"52":  "MX",
	"34":  "ES",
	"351": "PT",
	"44":  "GB",
	"49":  "DE",
	"33":  "FR",
	"39":  "IT",
	"55":  "BR",
	"54":  "AR",
	"81":  "JP",
	"86":  "CN",
	"91":  "IN",
	"61":  "AU",
	"64":  "NZ",
}

var dialingCodesByLenDesc = []string{
	"351",
	"33", "34", "39", "44", "49", "52", "54", "55", "61", "64", "81", "86", "91",
	"1",
}

var areaCodeLenByDialing = map[string]int{
	"1":   3, // (US/CA/Caribbean)
	"52":  3, // MX
	"34":  3, // ES
	"33":  1, // FR
	"55":  2, // BR
	"61":  1, // AU
	"64":  1, // NZ
	"351": 2, // PT
	// add more area code lengths to support more cases
	// 0 = ambiguous, requires user supplied spaces
	"44": 0, // GB
	"49": 0, // DE
	"39": 0, // IT
	"54": 0, // AR
	"81": 0, // JP
	"86": 0, // CN
	"91": 0, // IN
}

func dialingForISO(iso string) string {
	return isoToDialing[strings.ToUpper(iso)]
}

func lookupByDialing(digits string) (string, string) {
	for _, d := range dialingCodesByLenDesc {
		if strings.HasPrefix(digits, d) {
			return d, dialingToISO[d]
		}
	}
	return "", ""
}
