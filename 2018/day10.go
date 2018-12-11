package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

type Particle struct {
	PosX, PosY int
	VelX, VelY int
}

type Coord struct {
	X, Y int
}

func (p *Particle) Tick() Coord {
	p.PosX += p.VelX
	p.PosY += p.VelY
	return Coord{p.PosX, p.PosY}
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	particles := make(map[*Particle]bool)
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		posX, _ := strconv.Atoi(strings.TrimLeft(l[10:16], " "))
		posY, _ := strconv.Atoi(strings.TrimLeft(l[18:24], " "))
		velX, _ := strconv.Atoi(strings.TrimLeft(l[36:38], " "))
		velY, _ := strconv.Atoi(strings.TrimLeft(l[40:42], " "))
		p := Particle{posX, posY, velX, velY}
		particles[&p] = true
	}

	for secs := 1; ; secs++ {
		seenX := make(map[int]int)
		seenY := make(map[int]int)
		picture := make(map[Coord]bool)
		minX := 100000
		maxX := -100000
		minY := 100000
		maxY := -100000
		for k, _ := range particles {
			pos := k.Tick()
			seenX[pos.X] += 1
			seenY[pos.Y] += 1
			if minX > pos.X {
				minX = pos.X
			}
			if maxX < pos.X {
				maxX = pos.X
			}
			if minY > pos.Y {
				minY = pos.Y
			}
			if maxY < pos.Y {
				maxY = pos.Y
			}
			picture[pos] = true
		}
		candidateX := false
		candidateY := false
		for _, v := range seenX {
			if v > 20 {
				candidateX = true
			}
		}
		for _, v := range seenY {
			if v > 20 {
				candidateY = true
			}
		}
		if candidateX && candidateY {
			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {
					if picture[Coord{x, y}] {
						fmt.Printf("#")
					} else {
						fmt.Printf(".")
					}
				}
				fmt.Println()
			}
			fmt.Printf("Seen in %d secs\n", secs)
			break
		}
	}
}
