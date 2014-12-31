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
  value string
}

func drawUI(renderer *sdl.Renderer) {
  for k := range elements {
    surface := ubuntuR.RenderText_Solid(k.text+": "+k.value, sdl.Color{255, 255, 0, 255})

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

func addUiElement(pos pos, text string) *uiElement {
  element := &uiElement{pos, text, "-"}
  elements[element] = struct{}{}
  return element
}
