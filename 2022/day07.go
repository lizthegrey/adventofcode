package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

type dir struct {
	parent   *dir
	children []*dir
	size     int
}

func (d dir) walk(cb func(dir)) {
	cb(d)
	for _, v := range d.children {
		v.walk(cb)
	}
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var root dir
	cur := &root
	for _, s := range split[1 : len(split)-1] {
		tokens := strings.Split(s, " ")
		if tokens[0] == "$" {
			switch tokens[1] {
			case "cd":
				target := tokens[2]
				if target == ".." {
					cur = cur.parent
				} else {
					// Lazy-create the directory once we recurse into it.
					var child dir
					child.parent = cur
					cur.children = append(cur.children, &child)
					cur = &child
				}
			case "ls":
				// This will be handled by the else case.
				continue
			default:
				fmt.Printf("Bad command: %s\n", tokens[1])
				return
			}
		} else {
			// This is the result of listing a directory.
			if tokens[0] == "dir" {
				// Created when we recurse into it. Do nothing.
			} else {
				size, err := strconv.Atoi(tokens[0])
				if err != nil {
					fmt.Printf("Failed to parse size: %v\n", err)
					return
				}
				update := cur
				for update != nil {
					update.size += size
					update = update.parent
				}
			}
		}
	}

	// part A
	var total int
	root.walk(func(d dir) {
		if d.size <= 100000 {
			total += d.size
		}
	})
	fmt.Println(total)

	// part B
	var smallest int
	root.walk(func(d dir) {
		if root.size-d.size < 70000000-30000000 {
			if smallest == 0 || d.size < smallest {
				smallest = d.size
			}
		}
	})
	fmt.Println(smallest)
}
