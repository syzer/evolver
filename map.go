package main

import (
  // "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "math/rand"
  // "os"
)

type world struct {
  sections      map[pos]*section
  renderer      *sdl.Renderer
  height, width int64
  sectionsSize  int64
  sectionsCount int64
}

func createWorld(renderer *sdl.Renderer) (w world) {
  w.sections = make(map[pos]*section)
  w.sectionsCount = 20
  w.sectionsSize = 100
  for i := int64(0); i < w.sectionsCount; i++ {
    for j := int64(0); j < w.sectionsCount; j++ {
      section := section{
        pos:     pos{i, j},
        plants:  make(map[*plant]interface{}),
        animals: make(map[*animal]interface{})}
      w.sections[pos{i, j}] = &section
    }
  }
  w.renderer = renderer
  w.height = w.sectionsCount * w.sectionsSize
  w.width = w.sectionsCount * w.sectionsSize
  for i := 0; i < 1000; i++ {
    w.addRandomPlant()
  }

  for i := 0; i < 1000; i++ {
    w.addRandomAnimal()
  }

  return
}

func (w *world) addRandomPlant() {
  x := rand.Int63n(w.width)
  y := rand.Int63n(w.height)
  section := w.sections[pos{x / w.sectionsSize, y / w.sectionsSize}]
  p := plant{
    pos:     pos{x: x, y: y},
    section: section,
  }
  section.plants[&p] = struct{}{}
}

func (w *world) addRandomAnimal() {
  x := rand.Int63n(w.width)
  y := rand.Int63n(w.height)
  section := w.sections[pos{x / w.sectionsSize, y / w.sectionsSize}]
  a := animal{
    pos:     pos{x: x, y: y},
    section: section,
  }
  section.animals[&a] = struct{}{}
}

type section struct {
  pos
  plants  map[*plant]interface{}
  animals map[*animal]interface{}
}

type pos struct {
  x, y int64
}

type plant struct {
  pos
  section *section
}

type animal struct {
  pos
  section *section
}

func (m *world) makeTurn() {
  for _, section := range m.sections {
    // for p := range section.plants {
    // }

    for a := range section.animals {
      a.x += 1
    }
  }
}

func (m *world) draw() {

  for _, section := range m.sections {
    m.renderer.SetDrawColor(0, 255, 0, 255)
    for p := range section.plants {
      m.renderer.DrawPoint(int(p.x), int(p.y))
    }

    m.renderer.SetDrawColor(255, 0, 0, 255)
    for a := range section.animals {
      m.renderer.DrawPoint(int(a.x), int(a.y))
    }
  }

  // renderer.SetDrawColor(0, 0, 255, 255)
  // renderer.DrawLine(0, 0, 200, 200)

  points := []sdl.Point{{0, 0}, {100, 300}, {100, 300}, {200, 0}}
  m.renderer.SetDrawColor(255, 255, 0, 255)
  m.renderer.DrawPoints(points)

  // rect = sdl.Rect{300, 0, 200, 200}
  // renderer.SetDrawColor(255, 0, 0, 255)
  // renderer.DrawRect(&rect)

  // rects = []sdl.Rect{{400, 400, 100, 100}, {550, 350, 200, 200}}
  // renderer.SetDrawColor(0, 255, 255, 255)
  // renderer.DrawRects(rects)

  // rect = sdl.Rect{250, 250, 200, 200}
  // renderer.SetDrawColor(0, 255, 0, 255)
  // renderer.FillRect(&rect)

  // rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
  // renderer.SetDrawColor(255, 0, 255, 255)
  // renderer.FillRects(rects)
}
