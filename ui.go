package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  // "github.com/veandco/go-sdl2/sdl_ttf"
  "os"
  // "time"
)

var elements map[*uiElement]struct{} = make(map[*uiElement]struct{})

type uiElement struct {
  pos
  text  string
  fn    func() string
  color sdl.Color
}

func drawUI(renderer *sdl.Renderer) {
  for k := range elements {
    surface, _ := ubuntuR.RenderUTF8_Solid(k.text+": "+k.fn(), k.color)

    txt, err2 := renderer.CreateTextureFromSurface(surface)
    if err2 != nil {
      fmt.Fprintf(os.Stderr, "Failed to create texture from surface: %s\n", err2)
      os.Exit(2)
    }
    renderer.Copy(txt, nil, &sdl.Rect{int32(k.x), int32(k.y), surface.W, surface.H})

    surface.Free()
    txt.Destroy()
  }
}

func addWhiteUiElement(pos pos, text string, fn func() string) *uiElement {
  element := &uiElement{pos, text, fn, sdl.Color{255, 255, 255, 255}}
  elements[element] = struct{}{}
  return element
}

func addUiElement(pos pos, text string, fn func() string, color sdl.Color) *uiElement {
  element := &uiElement{pos, text, fn, color}
  elements[element] = struct{}{}
  return element
}
