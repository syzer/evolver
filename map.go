package main

import (
  // "fmt"
  "github.com/veandco/go-sdl2/sdl"
  // "os"
)

type world struct {
  sections map[pos]*section
  renderer *sdl.Renderer
}

func createWorld(renderer *sdl.Renderer) (w world) {
  w.sections = make(map[pos]*section)
  w.renderer = renderer
  return
}

type section struct {
  pos
}

type pos struct {
  x, y int32
}

func (m *world) draw() {
  m.renderer.SetDrawColor(0, 0, 0, 255)
  m.renderer.Clear()

  m.renderer.SetDrawColor(255, 255, 255, 255)
  m.renderer.DrawPoint(150, 300)

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

  m.renderer.Present()
}
