package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

type BitArray []bool

func (b BitArray) Read(start *int, length int) uint64 {
	if length > 64 || *start+length > len(b) {
		// explode
		fmt.Printf("should never get here: called with start %d and length %d\n", *start, length)
		return 0
	}
	var ret uint64
	for i := 0; i < length; i++ {
		ret = ret << 1
		if b[i+*start] {
			ret += 1
		}
	}
	*start += length
	return ret
}

func (b BitArray) Print() {
	for _, v := range b {
		if v {
			fmt.Printf("1")
		} else {
			fmt.Printf("0")
		}
	}
	fmt.Println()
}

type Packet struct {
	Version, Type int
	Value         *uint64
	Contained     []*Packet
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	input := make(BitArray, 0)
	for i := 0; i < len(split[0]); i += 2 {
		thisByte, _ := strconv.ParseInt(split[0][i:i+2], 16, 16)
		for n := 0; n < 8; n++ {
			input = append(input, (thisByte>>(7-n))&1 == 1)
		}
	}

	pos := 0
	outer := Parse(input, &pos)
	fmt.Println(outer.VersionSum())
	fmt.Println(outer.Evaluate())
}

func (p Packet) Evaluate() uint64 {
	switch p.Type {
	case 0:
		var sum uint64
		for _, inner := range p.Contained {
			sum += inner.Evaluate()
		}
		return sum
	case 1:
		product := uint64(1)
		for _, inner := range p.Contained {
			product *= inner.Evaluate()
		}
		return product
	case 2:
		min := uint64(math.MaxUint64)
		for _, inner := range p.Contained {
			val := inner.Evaluate()
			if val < min {
				min = val
			}
		}
		return min
	case 3:
		var max uint64
		for _, inner := range p.Contained {
			val := inner.Evaluate()
			if val > max {
				max = val
			}
		}
		return max
	case 4:
		return *p.Value
	case 5:
		if p.Contained[0].Evaluate() > p.Contained[1].Evaluate() {
			return uint64(1)
		}
		return uint64(0)
	case 6:
		if p.Contained[0].Evaluate() < p.Contained[1].Evaluate() {
			return uint64(1)
		}
		return uint64(0)
	case 7:
		if p.Contained[0].Evaluate() == p.Contained[1].Evaluate() {
			return uint64(1)
		}
		return uint64(0)
	default:
		fmt.Printf("Got invalid type %d\n", p.Type)
		return uint64(0)
	}
}

func (p Packet) VersionSum() int {
	sum := p.Version
	for _, inner := range p.Contained {
		sum += inner.VersionSum()
	}
	return sum
}

func Parse(input BitArray, pos *int) *Packet {
	version := int(input.Read(pos, 3))
	packetType := int(input.Read(pos, 3))
	if packetType == 4 {
		// literal
		var value uint64
		for input.Read(pos, 1) == 1 {
			bits4 := input.Read(pos, 4)
			value = value << 4
			value += bits4
		}
		bits4 := input.Read(pos, 4)
		value = value << 4
		value += bits4
		return &Packet{
			Version: version,
			Type:    packetType,
			Value:   &value,
		}
	}
	// operator. we will have children.
	var contained []*Packet
	if lengthType := input.Read(pos, 1); lengthType == 1 {
		subPacketCount := int(input.Read(pos, 11))
		for len(contained) < subPacketCount {
			contained = append(contained, Parse(input, pos))
		}
	} else {
		additionalBytes := int(input.Read(pos, 15))
		end := *pos + additionalBytes
		for *pos < end {
			contained = append(contained, Parse(input, pos))
		}
	}
	return &Packet{
		Version:   version,
		Type:      packetType,
		Contained: contained,
	}
}
