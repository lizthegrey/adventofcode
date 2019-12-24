package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"time"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Packet struct {
	X, Y int
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	var inputs, outputs [50]chan int
	var blocked [50]*bool
	done := make(chan bool)
	var nat Packet

	for i := 0; i < 50; i++ {
		input := make(chan int, 500)
		inputs[i] = input
		input <- i
		workingTape := tape.Copy()
		tid := -1
		if i == 0 {
			// tid = i
		}
		output, _, bl := workingTape.ProcessNonBlocking(input, tid)
		outputs[i] = output
		blocked[i] = bl
	}

	// Set up forwarder goroutines.
	for i := range outputs {
		go func(src int) {
			ch := outputs[src]
			for {
				dst := <-ch
				x := <-ch
				y := <-ch
				if dst == 255 {
					fmt.Printf("%d [%t %d]->255: %d,%d\n",
						src, *blocked[src], len(inputs[src]), x, y)
					if !*partB {
						fmt.Println(y)
						done <- true
					} else {
						nat = Packet{x, y}
					}
					continue
				}
				fmt.Printf("%d [%t %d]->%d [%t %d]: %d,%d\n",
					src, *blocked[src], len(inputs[src]),
					dst, *blocked[dst], len(inputs[dst]),
					x, y)
				inputs[dst] <- x
				inputs[dst] <- y
			}
		}(i)
	}

	go func() {
		var last Packet
		for {
			time.Sleep(20 * time.Millisecond)
			idle := true
			for i, bl := range blocked {
				if !*bl || len(inputs[i]) != 0 {
					idle = false
					break
				}
			}
			if idle {
				fmt.Printf("Detected idle, sending %d,%d to 0\n", nat.X, nat.Y)
				if last.Y == nat.Y {
					fmt.Printf("Saw repeated Y: %d\n", last.Y)
					done <- true
					return
				} else {
					last = nat
				}
				inputs[0] <- nat.X
				inputs[0] <- nat.Y
			}
		}
	}()

	<-done
}
