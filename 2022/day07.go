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
	children map[string]*dir
	leaves   map[string]int
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

	var cwd []string
	root := dir{
		parent:   nil,
		children: make(map[string]*dir),
		leaves:   make(map[string]int),
		size:     0,
	}
	cur := &root
	for _, s := range split[1 : len(split)-1] {
		tokens := strings.Split(s, " ")
		if tokens[0] == "$" {
			switch tokens[1] {
			case "cd":
				target := tokens[2]
				if target == ".." {
					cwd = cwd[:len(cwd)-1]
					cur = cur.parent
				} else {
					cwd = append(cwd, target)
					cur = cur.children[target]
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
				cur.children[tokens[1]] = &dir{
					parent:   cur,
					children: make(map[string]*dir),
					leaves:   make(map[string]int),
					size:     0,
				}
			} else {
				size, err := strconv.Atoi(tokens[0])
				if err != nil {
					fmt.Printf("Failed to parse size: %v\n", err)
					return
				}
				cur.leaves[tokens[1]] = size
				update := cur
				for i := 0; i <= len(cwd); i++ {
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
