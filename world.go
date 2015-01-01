package main

import (
  // "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "math"
  "math/rand"
  "sync"
  // "os"
)

var routines int = 4
var currentId int64 = 0

type world struct {
  sections      map[posInt]*section
  renderer      *sdl.Renderer
  height, width int64
  sectionsSize  int32
  sectionsCount int32
}

func createWorld(renderer *sdl.Renderer) (w world) {
  w.sections = make(map[posInt]*section)
  w.sectionsCount = 7
  w.sectionsSize = 100
  for i := int32(0); i < w.sectionsCount; i++ {
    for j := int32(0); j < w.sectionsCount; j++ {
      section := section{
        posInt:  posInt{i, j},
        plants:  make(map[*plant]interface{}),
        animals: make(map[*animal]interface{})}
      w.sections[posInt{i, j}] = &section
    }
  }
  w.renderer = renderer
  w.height = int64(w.sectionsCount) * int64(w.sectionsSize)
  w.width = int64(w.sectionsCount) * int64(w.sectionsSize)
  for i := 0; i < 10000; i++ {
    w.addRandomPlant()
  }

  for i := 0; i < 20; i++ {
    w.addRandomAnimal()
  }

  return
}

func (w *world) addRandomPlant() {
  x := float64(rand.Int63n(w.width))
  y := float64(rand.Int63n(w.height))
  secPos := posInt{int32(math.Floor(x / float64(w.sectionsSize))), int32(math.Floor(y / float64(w.sectionsSize)))}
  section := w.sections[secPos]
  p := plant{
    pos:     pos{x: x, y: y},
    section: section,
  }
  section.plants[&p] = struct{}{}
}

func (w *world) addRandomAnimal() {
  x := float64(rand.Int63n(w.width))
  y := float64(rand.Int63n(w.height))
  secPos := posInt{int32(math.Floor(x / float64(w.sectionsSize))), int32(math.Floor(y / float64(w.sectionsSize)))}
  section := w.sections[secPos]
  a := animal{
    pos:     pos{x: x, y: y},
    section: section,
    dMove:   nil,
    id:      currentId,
  }
  currentId++
  section.animals[&a] = struct{}{}
}

type section struct {
  posInt
  plants  map[*plant]interface{}
  animals map[*animal]interface{}
}

type posInt struct {
  x, y int32
}

type pos struct {
  x, y float64
}

type plant struct {
  pos
  section *section
}

type animal struct {
  pos
  section *section
  // decisions
  dMove     *pos
  dEatPlant *plant
  id        int64
}

func (w *world) makeTurn() {
  var wg sync.WaitGroup
  wg.Add(routines)
  for i := 0; i < routines; i++ {
    go w.makeDecision(i, &wg)
  }
  wg.Wait()
  for _, section := range w.sections {
    // for p := range section.plants {
    // }

    for a := range section.animals {
      if a.dMove != nil {
        a.x += a.dMove.x
        a.y += a.dMove.y
        // get new section
        secPos := posInt{int32(math.Floor(a.x / float64(w.sectionsSize))), int32(math.Floor(a.y / float64(w.sectionsSize)))}
        if secPos != a.section.posInt {
          // fmt.Println("Moving from", a.section.posInt, secPos, a.id)
          delete(a.section.animals, a)
          w.sections[secPos].animals[a] = struct{}{}
          a.section = w.sections[secPos]
        }
      }
      if a.dEatPlant != nil {
        plant := a.dEatPlant
        if plant.section != nil {
          delete(plant.section.plants, plant)
          // plant.section.plants[plant] = nil
          plant.section = nil
        }
      }

    }
  }
}

func (w *world) makeDecision(num int, wg *sync.WaitGroup) {
  for _, section := range w.sections {
    if int(section.y)%routines != num {
      // another thread will take care of this.
      continue
    }
    for a := range section.animals {
      w.animalAi(a)
    }
  }
  wg.Done()
}

var sightRange float64 = 50
var speed float64 = 1.5

func (w *world) animalAi(a *animal) {
  // watch out, is there some food around?
  a.dMove = nil
  a.dEatPlant = nil
  var closestPlant *plant = nil
  var plantDist = sightRange
  for i := a.section.x - 1; i <= a.section.x+1; i++ {
    for j := a.section.y - 1; j <= a.section.y+1; j++ {
      s := w.sections[posInt{i, j}]
      if s != nil {

        for p := range s.plants {
          if p.distance(&a.pos) <= plantDist {
            plantDist = p.distance(&a.pos)
            closestPlant = p
          }
        }
      }
    }

  }
  if closestPlant != nil {
    var xChange, yChange float64
    if speed >= plantDist {
      xChange = (closestPlant.x - a.x)
      yChange = (closestPlant.y - a.y)
      a.dEatPlant = closestPlant
    } else {
      xChange = (closestPlant.x - a.x) * (speed / plantDist)
      yChange = (closestPlant.y - a.y) * (speed / plantDist)
    }
    a.dMove = &pos{xChange, yChange}

  } else {
    // fmt.Println("No more food...", a.id)
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

func (p1 *pos) distance(p2 *pos) float64 {
  return math.Sqrt((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y))
}
