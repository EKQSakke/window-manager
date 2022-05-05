package main

import (
    "os"
    "strings"
    "strconv"
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

func (layout *LayoutWindow) toString(){
    println(layout.x, layout.y, layout.w, layout.h)
}

func getLayout(layoutId int) {
    f, err := os.ReadFile("layouts.txt")
    Check(err)
    input := string(f)

    newLayout := Layout{}
    layout := strings.Split(input, "\n")[layoutId]
    windows := strings.Split(layout, "; ")
    for _, window := range windows {
        props := strings.Split(window, ",")
        newLayout.AddWindowToLayout(props)
    }
    for _, w := range newLayout.Windows {
        w.toString()
    }
}
