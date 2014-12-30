package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/sdl_ttf"
  "os"
  "time"
)

var winTitle string = "Evolver"
var winWidth, winHeight int = 1200, 800

var running = true

var ubuntuR, ubuntuB *ttf.Font

func init() {
  if ttf.Init() != 0 {
    fmt.Fprintf(os.Stderr, "Failed to init ttf: %s\n")
    os.Exit(1)
  }

  var err error

  ubuntuR, err = ttf.OpenFont("UbuntuMono-R.ttf", 20)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to open regular font: %s\n", err)
    os.Exit(1)
  }
  ubuntuB, err = ttf.OpenFont("UbuntuMono-B.ttf", 20)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to open bold font: %s\n", err)
    os.Exit(1)
  }
}

func main() {

  var window *sdl.Window
  var renderer *sdl.Renderer
  // var rect sdl.Rect
  // var rects []sdl.Rect

  window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    winWidth, winHeight, sdl.WINDOW_SHOWN)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
    os.Exit(1)
  }

  renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
    os.Exit(2)
  }

  w := createWorld(renderer)
  last := time.Now()
  for running {
    since := time.Since(last)
    if since > time.Millisecond*100 {
      refresh(renderer, &w)
      last = time.Now()
    }

    time.Sleep(time.Millisecond * 10)

    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch t := event.(type) {
      case *sdl.QuitEvent:
        running = false
      case *sdl.MouseMotionEvent:
        // fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
        //   t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
      case *sdl.MouseButtonEvent:
        // fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
        //   t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
      case *sdl.MouseWheelEvent:
        // fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
        //   t.Timestamp, t.Type, t.Which, t.X, t.Y)
      case *sdl.KeyUpEvent:
        handleKey(t.Keysym.Sym)
        // fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
        //   t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
      }
    }

  }

  renderer.Destroy()
  window.Destroy()
  ttf.Quit()
}

func refresh(renderer *sdl.Renderer, w *world) {
  renderer.SetDrawColor(0, 0, 0, 255)
  renderer.Clear()
  w.draw()
  drawUI(renderer)
  renderer.Present()
}

func handleKey(code sdl.Keycode) {
  if code == sdl.K_ESCAPE {
    running = false
  }
}

func drawUI(renderer *sdl.Renderer) {
  surface := ubuntuR.RenderText_Solid("jeb z lasera pistoletem!!!", sdl.Color{255, 255, 0, 255})

  txt, err2 := renderer.CreateTextureFromSurface(surface)
  if err2 != nil {
    fmt.Fprintf(os.Stderr, "Failed to create texture from surface: %s\n", err2)
    os.Exit(2)
  }
  renderer.Copy(txt, nil, &sdl.Rect{0, 0, surface.W, surface.H})

  surface.Free()
  txt.Destroy()
}
