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
var winWidth, winHeight int = 1200, 800

var screenPos = pos{0, 0}

var running = true

var desiredFps int64 = 30

var ubuntuR, ubuntuB *ttf.Font

var fpsElement *uiElement
var fpsDesiredElement *uiElement
var turnElement *uiElement
var windowSize pos

func init() {
  if ttf.Init() != nil {
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
  width, height := window.GetSize()
  windowSize.x = float64(width)
  windowSize.y = float64(height)
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
  currentFramesCounter := 0

  w := createWorld(renderer)

  fpsElement = addWhiteUiElement(pos{5, 5}, "Current FPS", func() string {
    return strconv.Itoa(currentFramesCounter)
  })

  fpsDesiredElement = addWhiteUiElement(pos{5, 25}, "Desired FPS", func() string {
    return strconv.Itoa(int(desiredFps))
  })

  turnElement = addWhiteUiElement(pos{5, 45}, "Turn number", func() string {
    return strconv.Itoa(int(w.turnNumber))
  })

  addSubtypeUi(&w)
  last := time.Now()

  for running {

    since := time.Since(last)
    if since > time.Millisecond*time.Duration(1000/desiredFps) {
      // fmt.Println("Refreshing", last)
      framesCounter++
      if time.Now().Unix() != last.Unix() {
        // Update fps counter.
        currentFramesCounter = framesCounter
        framesCounter = 0
      }
      last = time.Now()
      w.makeTurn()
      w.makeTurn()

      refresh(renderer, &w)

    }
    time.Sleep(time.Millisecond * 9) // should be 10, but it is mroe accurate with 10

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
        if t.Keysym.Sym == sdl.K_ESCAPE {
          running = false
        }
      case *sdl.KeyDownEvent:
        handleKey(t.Keysym.Sym, &w)
      }

    }

  }

  renderer.Destroy()
  window.Destroy()
  ttf.Quit()
}

func addSubtypeUi(w *world) {
  for i, s := range w.subtypes {
    subtype := s
    addUiElement(
      pos{210 + float64(i*150), 5},
      fmt.Sprintf("%-7s", subtype.name),
      func() string {
        return strconv.Itoa(int(subtype.livingCount))
      },
      sdl.Color{subtype.red, subtype.green, subtype.blue, 255},
    )

    addUiElement(
      pos{210 + float64(i*150), 25},
      "Kills  ",
      func() string {
        return strconv.Itoa(int(subtype.killCount))
      },
      sdl.Color{subtype.red, subtype.green, subtype.blue, 255},
    )

    addUiElement(
      pos{210 + float64(i*150), 45},
      "Deaths ",
      func() string {
        return strconv.Itoa(int(subtype.killedCount))
      },
      sdl.Color{subtype.red, subtype.green, subtype.blue, 255},
    )

    addUiElement(
      pos{210 + float64(i*150), 65},
      "Starved",
      func() string {
        return strconv.Itoa(int(subtype.starvedCount))
      },
      sdl.Color{subtype.red, subtype.green, subtype.blue, 255},
    )

    addUiElement(
      pos{210 + float64(i*150), 85},
      "Old Age",
      func() string {
        return strconv.Itoa(int(subtype.oldAgeCount))
      },
      sdl.Color{subtype.red, subtype.green, subtype.blue, 255},
    )
  }

}

func refresh(renderer *sdl.Renderer, w *world) {
  // renderer.SetDrawColor(0, 0, uint8(frameId%255), 255)
  renderer.SetDrawColor(0, 0, 0, 255)
  renderer.Clear()

  w.draw(screenPos, windowSize)
  drawUI(renderer)

  renderer.Present()
}

func handleKey(code sdl.Keycode, w *world) {

  switch code {
  case sdl.K_ESCAPE:
    running = false
  case sdl.K_PLUS, sdl.K_EQUALS, sdl.K_KP_PLUS:
    if desiredFps < 60 {
      desiredFps++
    }
  case sdl.K_MINUS, sdl.K_KP_MINUS:
    if desiredFps > 1 {
      desiredFps--
    }

  case sdl.K_LEFT:
    screenPos.x -= 25
    if screenPos.x < 0 {
      screenPos.x = 0
    }
  case sdl.K_RIGHT:
    screenPos.x += 25
    if screenPos.x > float64(w.width)-windowSize.x {
      screenPos.x = float64(w.width) - windowSize.x
    }

  case sdl.K_UP:
    screenPos.y -= 50
    if screenPos.y < 0 {
      screenPos.y = 0
    }

  case sdl.K_DOWN:
    screenPos.y += 50
    if screenPos.y > float64(w.height)-windowSize.y {
      screenPos.y = float64(w.height) - windowSize.y
    }
  }
}
