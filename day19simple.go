package main

import (
	"fmt"
	"regexp"
)

func main() {
	substitutions := map[string][]string{
		"Al": {"ThF", "ThRnFAr"},
		"B":  {"BCa", "TiB", "TiRnFAr"},
		"Ca": {"CaCa", "PB", "PRnFAr", "SiRnFYFAr", "SiRnMgAr", "SiTh"},
		"F":  {"CaF", "PMg", "SiAl"},
		"H":  {"CRnAlAr", "CRnFYFYFAr", "CRnFYMgAr", "CRnMgYFAr", "HCa", "NRnFYFAr", "NRnMgAr", "NTh", "OB", "ORnFAr"},
		"Mg": {"BF", "TiMg"},
		"N":  {"CRnFAr", "HSi"},
		"O":  {"CRnFYFAr", "CRnMgAr", "HP", "NRnFAr", "OTi"},
		"P":  {"CaP", "PTi", "SiRnFAr"},
		"Si": {"CaSi"},
		"Th": {"ThCa"},
		"Ti": {"BP", "TiTi"},
		"e":  {"HF", "NAl", "OMg"},
	}

	medicine := "CRnCaSiRnBSiRnFArTiBPTiTiBFArPBCaSiThSiRnTiBPBPMgArCaSiRnTiMgArCaSiThCaSiRnFArRnSiRnFArTiTiBFArCaCaSiRnSiThCaCaSiRnMgArFYSiRnFYCaFArSiThCaSiThPBPTiMgArCaPRnSiAlArPBCaCaSiRnFYSiThCaRnFArArCaCaSiRnPBSiRnFArMgYCaCaCaCaSiThCaCaSiAlArCaCaSiRnPBSiAlArBCaCaCaCaSiThCaPBSiThPBPBCaSiRnFYFArSiThCaSiRnFArBCaCaSiRnFYFArSiThCaPBSiThCaSiRnPMgArRnFArPTiBCaPRnFArCaCaCaCaSiRnCaCaSiRnFYFArFArBCaSiThFArThSiThSiRnTiRnPMgArFArCaSiThCaPBCaSiRnBFArCaCaPRnCaCaPMgArSiRnFYFArCaSiThRnPBPMgAr"

	//medicine = "HOHOHO"
	//substitutions = map[string][]string {
	//	"e": {"H", "O"},
	//	"H": {"HO", "OH"},
	//	"O": {"HH"},
	//}

	// part (a)
	calibration := make(map[string]bool)
	iterate(medicine, substitutions, calibration)
	fmt.Println(len(calibration))

	// part (b)
	starter := "e"
	molecule := medicine
	iterationCount := 0

	for molecule != starter {
		result := reverseIterate(molecule, substitutions)
		if result == nil {
			fmt.Println("Fatal: no matches.")
			return
		}
		molecule = *result
		fmt.Println(len(molecule))
		iterationCount++
	}

	fmt.Println(iterationCount)
}

func iterate(input string, substitutions map[string][]string, outputs map[string]bool) {
	for k, v := range substitutions {
		r := regexp.MustCompile(k)
		matches := r.FindAllStringIndex(input, -1)
		for m := range matches {
			for i := range v {
				out := make([]byte, len(input)-len(k)+len(v[i]))
				copy(out[0:matches[m][0]], input[0:matches[m][0]])
				copy(out[matches[m][0]:matches[m][0]+len(v[i])], v[i][:])
				copy(out[matches[m][0]+len(v[i]):len(out)], input[matches[m][1]:len(input)])
				outputs[string(out)] = true
			}
		}
	}
}

func reverseIterate(input string, substitutions map[string][]string) *string {
	for k, v := range substitutions {
		for i := range v {
			r := regexp.MustCompile(v[i])
			matches := r.FindAllStringIndex(input, -1)
			for m := range matches {
				out := make([]byte, len(input)+len(k)-len(v[i]))
				copy(out[0:matches[m][0]], input[0:matches[m][0]])
				copy(out[matches[m][0]:matches[m][0]+len(k)], k[:])
				copy(out[matches[m][0]+len(k):len(out)], input[matches[m][1]:len(input)])
				result := string(out)
				return &result
			}
		}
	}
	return nil
}
