package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")
var hexColor = regexp.MustCompile("^#[0-9a-f]{6}$")
var pidDigits = regexp.MustCompile("^[0-9]{9}$")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	passports := []map[string]string{
		make(map[string]string),
	}
	i := 0
	for _, s := range split {
		if s == "" {
			i++
			passports = append(passports, make(map[string]string))
			continue
		}
		for _, kv := range strings.Split(s, " ") {
			kvParts := strings.Split(kv, ":")
			if len(kvParts) != 2 {
				fmt.Printf("Failed to parse %s\n", kv)
				return
			}
			passports[i][kvParts[0]] = kvParts[1]
		}
	}
	valid := 0
outer:
	for _, passport := range passports {
		for _, k := range []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"} {
			if passport[k] == "" {
				continue outer
			}
		}
		valid++
	}
	fmt.Println(valid)

	strict := 0
	rejects := make(map[string]int)
	for _, passport := range passports {
		byr, err := strconv.Atoi(passport["byr"])
		if err != nil || byr < 1920 || byr > 2002 {
			rejects["byr"]++
			continue
		}
		iyr, err := strconv.Atoi(passport["iyr"])
		if err != nil || iyr < 2010 || iyr > 2020 {
			rejects["iyr"]++
			continue
		}
		eyr, err := strconv.Atoi(passport["eyr"])
		if err != nil || eyr < 2020 || eyr > 2030 {
			rejects["eyr"]++
			continue
		}
		if len(passport["hgt"]) < 3 {
			rejects["hgt"]++
			continue
		}
		hgt, err := strconv.Atoi(passport["hgt"][:len(passport["hgt"])-2])
		units := passport["hgt"][len(passport["hgt"])-2:]
		if err != nil {
			rejects["hgt"]++
			continue
		}
		switch units {
		case "cm":
			if hgt < 150 || hgt > 193 {
				rejects["hgt"]++
				continue
			}
		case "in":
			if hgt < 59 || hgt > 76 {
				rejects["hgt"]++
				continue
			}
		default:
			rejects["hgt"]++
			continue
		}
		if !hexColor.MatchString(passport["hcl"]) {
			rejects["hcl"]++
			continue
		}
		switch passport["ecl"] {
		case "amb":
		case "blu":
		case "brn":
		case "gry":
		case "grn":
		case "hzl":
		case "oth":
		default:
			rejects["ecl"]++
			continue
		}
		if !pidDigits.MatchString(passport["pid"]) {
			rejects["pid"]++
			continue
		}
		strict++
	}
	fmt.Println(rejects)
	fmt.Println(strict)
}
