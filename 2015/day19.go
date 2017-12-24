package main

import (
	"fmt"
	"regexp"
	"sort"
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

	// Precompute some shortenings of common patterns
	for round := 0; round < 0; round++ {
		for k, v := range substitutions {
			newSubs := make(map[string]bool)
			for i := range v {
				iterate(v[i], substitutions, newSubs)
			}
			for ns := range newSubs {
				substitutions[k] = append(substitutions[k], ns)
			}
		}
	}

	known := map[int][]string{
		len(medicine): {medicine},
	}
	keys := []int{len(medicine)}
	iterationCount := 0
outer:
	for {
		if k, ok := known[len(starter)]; ok {
			for i := range k {
				if known[len(starter)][i] == starter {
					break outer
				}
			}
		}

		newKnown := make(map[int][]string)
		newKeys := make([]int, 0)
		for m := range keys {
			for v := range known[keys[m]] {
				newKeys = reverseIterate(known[keys[m]][v], substitutions, newKnown, newKeys)
			}
		}
		known = map[int][]string{
			newKeys[0]: newKnown[newKeys[0]],
		}
		keys = newKeys
		iterationCount++
		fmt.Printf("Shortest string is %d; we know about %d molecules\n", keys[0], len(known[keys[0]]))
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

func reverseIterate(input string, substitutions map[string][]string, outputs map[int][]string, keys []int) []int {
	for k, v := range substitutions {
		for i := range v {
			r := regexp.MustCompile(v[i])
			matches := r.FindAllStringIndex(input, -1)
			for m := range matches {
				out := make([]byte, len(input)+len(k)-len(v[i]))
				copy(out[0:matches[m][0]], input[0:matches[m][0]])
				copy(out[matches[m][0]:matches[m][0]+len(k)], k[:])
				copy(out[matches[m][0]+len(k):len(out)], input[matches[m][1]:len(input)])
				if outputs[len(out)] == nil {
					outputs[len(out)] = make([]string, 0)
					keys = append(keys, len(out))
				}
				outputs[len(out)] = append(outputs[len(out)], string(out))
			}
		}
	}
	sort.Ints(keys)
	return keys
}
