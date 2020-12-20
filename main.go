package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const sectionCount = 60
const bridgesCount = sectionCount / 5

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
	exitsCountProbabilities := [5]float32{0.1, 0.25, 0.80, 0.95, 1.0}
	n := rand.Float32()
	i := 0
	for n >= exitsCountProbabilities[i] {
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

func main() {
	sections := []*Section{}
	root := NewSection("root", 0)
	for root.MaxExits == 0 {
		root.MaxExits = pickMaxExit()
	}
	sections = append(sections, root)

	for i := 0; i < sectionCount; i++ {
		sec := NewSection(fmt.Sprintf("section #%v", i+1), i+1)

		for _, parent := range sections {
			if len(parent.Exits) < parent.MaxExits && parent.Depth < sectionCount/3 {
				parent.Exits = append(parent.Exits, sec)
				sec.Depth = parent.Depth + 1
				sec.Parent = parent
				sections = append(sections, sec)
				break
			}
		}
		if sec.Parent == nil {
			n := rand.Intn(len(sections))
			sections[n].MaxExits++
			sections[n].Exits = append(sections[n].Exits, sec)
			sections = append(sections, sec)
		}
	}

	for i, val := range rand.Perm(len(sections) - 1) {
		sections[i+1].Index = val + 1
	}

	for i := 0; i < bridgesCount; i++ {
		sec := sections[rand.Intn(len(sections))]
		bridge := sec
		for bridge == sec || absint(bridge.Depth-sec.Depth) > 2 {
			bridge = sections[rand.Intn(len(sections))]
		}
		sec.Exits = append(sec.Exits, bridge)
	}

	fmt.Println("digraph sections {\nnode [style=filled]\n0 [color=green]")
	for _, sec := range sections {
		if len(sec.Exits) == 0 {
			fmt.Printf("%v [color=blue]\n", sec.Index)
		}
		for _, exit := range sec.Exits {
			fmt.Printf("%s%v -> %v\n", strings.Repeat(" ", sec.Depth), sec.Index, exit.Index)
		}
	}
	fmt.Println("}")
}
