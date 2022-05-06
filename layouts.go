package main

import (
	"os"
	"strconv"
	"strings"
)

type Layout struct {
	Windows []LayoutWindow
}

func (layout *Layout) AddWindowToLayout(props []string) {
	x, _ := strconv.Atoi(props[0])
	y, _ := strconv.Atoi(props[1])
	w, _ := strconv.Atoi(props[2])
	h, _ := strconv.Atoi(props[3])
	layout.Windows = append(layout.Windows, LayoutWindow{x: x, y: y, w: w, h: h})
}

type LayoutWindow struct {
	x int
	y int
	w int
	h int
}

func (layout *LayoutWindow) toString() {
	println("x", layout.x, "y", layout.y, "w", layout.w, "h", layout.h)
}

func getLayout(layoutId int) Layout {
	println("layoutId:", layoutId)
	f, err := os.ReadFile("layouts.txt")
	Check(err)
	input := string(f)

	newLayout := Layout{}
	layout := strings.Split(input, "\n")[layoutId]
	windows := strings.Split(layout, ";")
	for _, window := range windows {
		props := strings.Split(window, ",")
        // last line probs
        if (len(props) < 4) {
            continue
        }
        println(props[0])
        println(props[1])
        println(props[2])
        println(props[3])
		newLayout.AddWindowToLayout(props)
	}
	for _, w := range newLayout.Windows {
		w.toString()
	}

	return newLayout
}
