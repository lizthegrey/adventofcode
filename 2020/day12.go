package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	// Normal 2d coord plane > x, ^ y
	posX := 0
	posY := 0

	facing := 0 // east
	equivalentFace := []byte{'E', 'N', 'W', 'S'}
	for _, s := range split {
		param, err := strconv.Atoi(s[1:len(s)])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}

		direction := s[0]
		switch s[0] {
		case 'R':
			facing -= (param / 90)
			facing += 4
			facing %= 4
		case 'L':
			facing += (param / 90)
			facing %= 4
		case 'F':
			direction = equivalentFace[facing]
		}

		switch direction {
		case 'N':
			posY += param
		case 'S':
			posY -= param
		case 'E':
			posX += param
		case 'W':
			posX -= param
		}
	}

	// We've reached the end.
	manhattan := 0
	if posX >= 0 {
		manhattan += posX
	} else {
		manhattan -= posX
	}
	if posY >= 0 {
		manhattan += posY
	} else {
		manhattan -= posY
	}
	fmt.Println(manhattan)

	// Part B.
	wayX := 10
	wayY := 1
	posX = 0
	posY = 0
	// wayX = A*wayX + B*wayY
	// wayY = C*wayX + D*wayY
	// {A,B,C,D}
	refl := [][4]int{
		// 0 degrees: (x,y)->(x,y)
		{1, 0, 0, 1},
		// 90 degrees: (x,y)->(-y,x)
		{0, -1, 1, 0},
		// 180 degrees: (x,y)->(-x,-y)
		{-1, 0, 0, -1},
		// 270 degrees: (x,y)->(y,-x)
		{0, 1, -1, 0},
	}

	for _, s := range split {
		param, err := strconv.Atoi(s[1:len(s)])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}

		switch s[0] {
		case 'R':
			idx := (4 + (-param / 90)) % 4
			newWayX := refl[idx][0]*wayX + refl[idx][1]*wayY
			newWayY := refl[idx][2]*wayX + refl[idx][3]*wayY
			wayX = newWayX
			wayY = newWayY
		case 'L':
			idx := (param / 90) % 4
			newWayX := refl[idx][0]*wayX + refl[idx][1]*wayY
			newWayY := refl[idx][2]*wayX + refl[idx][3]*wayY
			wayX = newWayX
			wayY = newWayY
		case 'F':
			posX += wayX * param
			posY += wayY * param
		case 'N':
			wayY += param
		case 'S':
			wayY -= param
		case 'E':
			wayX += param
		case 'W':
			wayX -= param
		}
	}

	// We've reached the end.
	manhattan = 0
	if posX >= 0 {
		manhattan += posX
	} else {
		manhattan -= posX
	}
	if posY >= 0 {
		manhattan += posY
	} else {
		manhattan -= posY
	}
	fmt.Println(manhattan)
}
