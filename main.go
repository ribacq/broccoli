package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const sectionCount = 60

var magic = int(math.Sqrt(float64(sectionCount)))

// Section of the interactive story
type Section struct {
	Index    int
	Title    string
	Text     string
	Depth    int
	MaxExits int
	Exits    []*Section
	Parent   *Section
}

// NewSection creates a Section with a title and index, and a randomly picked MaxExits
func NewSection(title string, index int) *Section {
	return &Section{
		Title:    title,
		Index:    index,
		MaxExits: pickMaxExit(),
	}
}

func pickMaxExit() int {
	exitsCountProbabilities := [5]float32{0.2, 0.4, 0.4, 0.1, 0.1}
	n := rand.Float32()
	i := 0
	for s := float32(0.0); n >= s; s += exitsCountProbabilities[i] {
		i++
	}
	return i
}

func absint(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Widths maps tree depths to widths
type Widths map[int]int

// NewWidths makes a new Widths map
func NewWidths() Widths {
	return make(map[int]int)
}

// Inc increments a level
func (w Widths) Inc(level int) {
	if _, ok := w[level]; !ok {
		w[level] = 0
	}
	w[level]++
}

// Dec decrements a level
func (w Widths) Dec(level int) {
	if _, ok := w[level]; !ok {
		w[level] = 0
	} else {
		w[level]--
	}
}

// Get returns the width at a given depth
func (w Widths) Get(level int) int {
	if _, ok := w[level]; !ok {
		w[level] = 0
	}
	return w[level]
}

// MaxDepth returns the maximum depth an the associated width
func (w Widths) MaxDepth() (maxDepth, deepWidth int) {
	for depth, width := range w {
		if depth > maxDepth {
			maxDepth = depth
			deepWidth = width
		}
	}
	return
}

// MaxWidth returns the maximum width and the associated depth
func (w Widths) MaxWidth() (maxWidth, wideDepth int) {
	for depth, width := range w {
		if width > maxWidth {
			maxWidth = width
			wideDepth = depth
		}
	}
	return
}

func main() {
	sections := []*Section{}
	root := NewSection("root", 0)
	for root.MaxExits == 0 {
		root.MaxExits = pickMaxExit()
	}
	sections = append(sections, root)

	widths := NewWidths()

	for i := 0; i < sectionCount; i++ {
		sec := NewSection(fmt.Sprintf("section #%v", i+1), i+1)

		for _, parent := range sections {
			if len(parent.Exits) < parent.MaxExits && widths.Get(parent.Depth+1) < magic {
				parent.Exits = append(parent.Exits, sec)
				sec.Depth = parent.Depth + 1
				widths.Inc(sec.Depth)
				sec.Parent = parent
				sections = append(sections, sec)
				break
			}
		}
		if sec.Parent == nil {
			n := rand.Intn(len(sections))
			sections[n].MaxExits++
			sections[n].Exits = append(sections[n].Exits, sec)
			sec.Depth = sections[n].Depth + 1
			widths.Inc(sections[n].Depth + 1)
			sec.Parent = sections[n]
			sections = append(sections, sec)
		}
	}

	for i, val := range rand.Perm(len(sections) - 1) {
		sections[i+1].Index = val + 1
	}

	for i := 0; i < magic; i++ {
		sec := sections[rand.Intn(len(sections))]
		bridge := sec
		for bridge == sec || absint(bridge.Depth-sec.Depth) > 2 {
			bridge = sections[rand.Intn(len(sections))]
		}
		sec.Exits = append(sec.Exits, bridge)
		if bridge.Depth > sec.Depth+1 {
			widths.Dec(bridge.Depth)
			bridge.Depth = sec.Depth + 1
			widths.Inc(bridge.Depth)
		}
	}

	fmt.Println("digraph sections {\nnode [style=filled]\n0 [color=green]")
	for _, sec := range sections {
		fmt.Printf("%v [label=\"%v\"", sec.Index, sec.Index)
		if len(sec.Exits) == 0 {
			fmt.Printf(",color=blue")
		}
		fmt.Println("]")
		for _, exit := range sec.Exits {
			fmt.Printf("%s%v -> %v\n", strings.Repeat(" ", sec.Depth), sec.Index, exit.Index)
		}
	}
	fmt.Println("}")

	md, dw := widths.MaxDepth()
	mw, wd := widths.MaxWidth()
	log.Printf("w.MaxDepth() = %v, %v\n", md, dw)
	log.Printf("w.MaxWidth() = %v, %v\n", mw, wd)
}
