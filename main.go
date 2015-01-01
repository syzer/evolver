package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/sdl_ttf"
  "os"
  "strconv"
  "time"
)

var winTitle string = "Evolver"
var winWidth, winHeight int = 800, 800

var running = true

var desiredFps int64 = 30

var ubuntuR, ubuntuB *ttf.Font

var fpsElement *uiElement
var fpsDesiredElement *uiElement

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

  fpsElement = addUiElement(pos{5, 5}, "Current FPS")
  fpsDesiredElement = addUiElement(pos{5, 25}, "Desired FPS")
  fpsDesiredElement.value = strconv.FormatInt(desiredFps, 10)
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

  renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
    os.Exit(2)
  }

  framesCounter := 0

  w := createWorld(renderer)
  last := time.Now()

  for running {

    since := time.Since(last)
    if since > time.Millisecond*time.Duration(1000/desiredFps) {
      // fmt.Println("Refreshing", last)
      framesCounter++
      if time.Now().Unix() != last.Unix() {
        // Update fps counter.
        fpsElement.value = strconv.Itoa(framesCounter)
        framesCounter = 0
      }
      last = time.Now()
      w.makeTurn()
      w.makeTurn()
      refresh(renderer, &w)

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
  // renderer.SetDrawColor(0, 0, uint8(frameId%255), 255)
  renderer.SetDrawColor(0, 0, 0, 255)
  renderer.Clear()
  w.draw()
  drawUI(renderer)

  renderer.Present()
}

func handleKey(code sdl.Keycode) {

  switch code {
  case sdl.K_ESCAPE:
    running = false
  case sdl.K_PLUS, sdl.K_EQUALS:
    if desiredFps < 60 {
      desiredFps++
      fpsDesiredElement.value = strconv.FormatInt(desiredFps, 10)
    }
  case sdl.K_MINUS:
    if desiredFps > 1 {
      desiredFps--
      fpsDesiredElement.value = strconv.FormatInt(desiredFps, 10)
    }
  }
}
