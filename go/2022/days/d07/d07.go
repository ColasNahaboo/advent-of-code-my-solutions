// Adventofcode 2022, d07, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 95437
// TEST: -1 input
// TEST: example 24933642
// TEST: input 5883165
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

var verbose bool

type Dir struct {
	name string
	parent *Dir					// nil at root
	dirs map[string]*Dir
	files map[string]int
	size int					// sum of all the files inside it
}

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) int {
	root := Dir{name: "/", parent: nil, dirs: make(map[string]*Dir, 0), files: make(map[string]int, 0)}
	parse(&root, lines, 0)
	du(&root)
	if verbose { VPdir(&root, "-"); }
	return duLess(&root, 100000)
}

//////////// Part 2
func part2(lines []string) int {
	root := Dir{name: "/", parent: nil, dirs: make(map[string]*Dir, 0), files: make(map[string]int, 0)}
	parse(&root, lines, 0)
	du(&root)
	if verbose { VPdir(&root, "-"); }
	return delFor(&root, 30000000 + root.size - 70000000, root.size)
}

//////////// Common Parts code

var recd = regexp.MustCompile("^[$] cd[[:space:]]+([-.[:alnum:]]+)")
var recommand = regexp.MustCompile("^[$]")
var refile = regexp.MustCompile("^([[:digit:]]+) ([-.[:alnum:]]+)")
var redir = regexp.MustCompile("^dir ([-.[:alnum:]]+)")

func parse(wd *Dir, lines []string, lineno int) *Dir {
	for i := 0 ; i < len(lines); i++ {
		line := lines[i]
		if line == "$ cd /" {
			for wd.parent != nil {
				wd = wd.parent
			}
			continue
		} else if line == "$ cd .." {
			if wd.parent == nil {
				log.Fatalln("cd .. at root!, at line", i+1)
			}
			wd = wd.parent
			continue
		} else if rescd := recd.FindStringSubmatch(line); rescd != nil {
			dirname := rescd[1]
			if wd.dirs[dirname] == nil {
				wd.dirs[dirname] = &Dir{dirname, wd, make(map[string]*Dir, 0), make(map[string]int, 0), 0}
			}
			wd = wd.dirs[dirname]
		} else if line == "$ ls" {
			for ;; {
				if (i + 1) >= len(lines) || recommand.MatchString(lines[i+1]) {
					break
				}
				i++
				line = lines[i]
				if resfile := refile.FindStringSubmatch(line); resfile != nil {
					filename := resfile[2]
					wd.files[filename] = atoi(resfile[1])
				} else if resdir := redir.FindStringSubmatch(line); resdir != nil {
					dirname := resdir[1]
					if wd.dirs[dirname] == nil {
						wd.dirs[dirname] = &Dir{dirname, wd, make(map[string]*Dir, 0), make(map[string]int, 0), 0}
					}
				} else {
					log.Fatalln("syntax error in ls output, at line", i+1, ":", line)
				}
			}
		} else {
			log.Fatalln("Syntax error at line", i+1, ":", line)
		}
	}
	return wd
}

// recursively computes all dir sizes and cache them in tree
func du(d *Dir) int {
	if d.size != 0 {
		return d.size
	}
	size := 0
	for _, sd := range d.dirs {
		size += du(sd)
	}
	for _, fsize := range d.files {
		size += fsize
	}
	d.size = size
	return size
}

// pretty-print the tree
func VPdir(d *Dir, indent string) {
	fmt.Printf("  %s %s (dir,  size=%d)\n", indent, d.name, d.size)
	for _, sd := range d.dirs {
		VPdir(sd, "  " + indent)
	}
	for name, size := range d.files {
		fmt.Printf("  %s %s (file, size=%d)\n", indent, name, size)
	}
}
	

//////////// Part1 functions

// sum of all dirs less than max (100k) each
// du must have been run before

func duLess(d *Dir, max int) int {
	size := 0
	for _, sd := range d.dirs {
		size += duLess(sd, max)
	}
	if d.size <= max {
		size += d.size
	}
	return size
}

//////////// Part2 functions

// find the size of the smallest dir to delete to free space space
// du must have been run before

func delFor(d *Dir, space, found int) int {
	
	for _, sd := range d.dirs {
		if newfound := delFor(sd, space, found); newfound > 0 && newfound < found {
			found = newfound
		}
	}
	if d.size >= space && d.size < found {
		found = d.size
	}
	return found
}
