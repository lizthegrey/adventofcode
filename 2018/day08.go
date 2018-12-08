package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")
var verbose = flag.Bool("verbose", false, "Print verbose debugging output.")

type Node struct {
	Parent   *Node
	Children []*Node
	Metadata []int
	Value    int
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	l, err := r.ReadString('\n')
	if err != nil || len(l) == 0 {
		return
	}
	l = l[:len(l)-1]
	entries := strings.Split(l, " ")
	values := make([]int, len(entries))
	for i, n := range entries {
		v, _ := strconv.Atoi(n)
		values[i] = v
	}

	cursor := 0
	var root *Node
	var parent *Node
	runningTotalMetadata := 0
	for cursor != len(values) {
		// Assume we start at the beginning of an entry.
		childCount := values[cursor]
		cursor += 1
		metadataCount := values[cursor]
		cursor += 1

		if *verbose {
			fmt.Printf("Starting new node with %d children and %d metadata.\n", childCount, metadataCount)
		}

		currentNode := &Node{parent, make([]*Node, childCount), make([]int, metadataCount), 0}

		if childCount == 0 {
			for i := 0; i < len(currentNode.Metadata); i++ {
				currentNode.Metadata[i] = values[cursor]
				runningTotalMetadata += values[cursor]
				currentNode.Value += values[cursor]
				cursor += 1
			}
		}

		if parent == nil {
			root = currentNode
			if childCount != 0 {
				parent = currentNode
			}
			continue
		}

		siblingCount := 0
		for _, v := range currentNode.Parent.Children {
			if v != nil {
				siblingCount += 1
			}
		}
		currentNode.Parent.Children[siblingCount] = currentNode

		if childCount != 0 {
			parent = currentNode
			continue
		}

		for currentNode.Parent != nil {
			siblingCount = 0
			for _, v := range currentNode.Parent.Children {
				if v != nil {
					siblingCount += 1
				}
			}
			if siblingCount != len(currentNode.Parent.Children) {
				break
			}
			currentNode = currentNode.Parent
			parent = currentNode.Parent
			for i := 0; i < len(currentNode.Metadata); i++ {
				currentNode.Metadata[i] = values[cursor]
				runningTotalMetadata += values[cursor]
				index := values[cursor]
				if index != 0 && index <= len(currentNode.Children) {
					currentNode.Value += currentNode.Children[index-1].Value
				}
				cursor += 1
			}
		}
	}

	fmt.Printf("Part A result is %d\n", runningTotalMetadata)
	fmt.Printf("Part B result is %d\n", root.Value)
}
