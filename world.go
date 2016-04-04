package main

import (
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "math"
  "math/rand"
  "sync"
  // "os"
)

var routines int = 4
var currentId int64 = 0

type world struct {
  subtypes         []*subtype
  sections         map[posInt]*section
  renderer         *sdl.Renderer
  height, width    int64
  sectionsSize     int32
  sectionsCount    int32
  turnNumber       int32
  carnivouresCount int32
  herbivoresCount  int32
}

type subtype struct {
  name  string
  speed float64
  red   uint8
  green uint8
  blue  uint8

  strength int32
  hp       int32

  // stats
  livingCount  int32
  killCount    int32
  killedCount  int32
  starvedCount int32
  oldAgeCount  int32
}

func createWorld(renderer *sdl.Renderer) (w world) {
  w.createSubtypes()
  w.sections = make(map[posInt]*section)
  w.sectionsCount = 50
  w.sectionsSize = 50
  w.turnNumber = 0
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

  for i := 0; i < 200; i++ {
    w.addRandomAnimal()
  }

  return
}

func (w *world) addPlant(x float64, y float64) {
  secPos := posInt{int32(math.Floor(x / float64(w.sectionsSize))), int32(math.Floor(y / float64(w.sectionsSize)))}
  section := w.sections[secPos]
  if section == nil {
    return
  }
  p := plant{
    pos:     pos{x: x, y: y},
    section: section,
    age:     0,
  }
  section.plants[&p] = struct{}{}
}

func (w *world) addRandomPlant() {
  x := float64(rand.Int63n(w.width))
  y := float64(rand.Int63n(w.height))
  w.addPlant(x, y)
}

func (w *world) addRandomAnimal() {
  x := float64(rand.Int63n(w.width))
  y := float64(rand.Int63n(w.height))
  secPos := posInt{int32(math.Floor(x / float64(w.sectionsSize))), int32(math.Floor(y / float64(w.sectionsSize)))}
  section := w.sections[secPos]
  a := animal{
    pos:     pos{x: x, y: y},
    section: section,
    id:      currentId,
    food:    300,
    subtype: w.subtypes[rand.Int31n(int32(len(w.subtypes)))],
  }
  a.hp = a.subtype.hp
  a.red = a.subtype.red
  a.green = a.subtype.green
  a.blue = a.subtype.blue
  a.subtype.livingCount++
  currentId++
  section.animals[&a] = struct{}{}
}

func (w *world) birth(a *animal) {
  child := animal{
    pos:     pos{x: a.x, y: a.y},
    section: a.section,
    id:      currentId,
    subtype: a.subtype,
    hp:      a.subtype.hp,
    food:    200,
    red:     a.red,
    green:   a.green,
    blue:    a.blue,
  }
  currentId++
  a.subtype.livingCount++
  a.section.animals[&child] = struct{}{}
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
  birth   bool
  dead    bool
  age     int32
}

type animal struct {
  pos
  subtype *subtype
  section *section
  dead    bool
  birth   bool
  id      int64
  food    int32
  age     int32
  hp      int32
  // decisions
  dMove     *pos
  dEatPlant *plant
  dAttack   *animal
  // ai
  wandering *pos

  // fenotype
  red uint8
  green uint8
  blue uint8
  gen gen
}

func (w *world) makeTurn() {
  w.turnNumber++
  w.decisionsPhase()
  // change world after the turn. Grow plants.
  w.executionPhase()
  w.mapModifyPhase()

}

func (w *world) decisionsPhase() {
  var wg sync.WaitGroup
  wg.Add(routines)
  for i := 0; i < routines; i++ {
    go w.makeDecision(i, &wg)
  }
  wg.Wait()
}

func (w *world) executionPhase() {
  for _, section := range w.sections {
    // for p := range section.plants {
    // }

    for a := range section.animals {
      if a.dMove != nil {
        a.x += a.dMove.x
        a.y += a.dMove.y
      }
      if a.dEatPlant != nil {
        plant := a.dEatPlant
        if !plant.dead {
          plant.dead = true
          a.food += 100
        }
      }
      if a.dAttack != nil && !a.dAttack.dead {
        damage := rand.Int31n(a.subtype.strength)
        if a.dAttack.hp <= damage {
          a.dAttack.dead = true
          a.dAttack.subtype.killedCount++
          a.subtype.killCount++
          a.food += 500
        } else {
          a.dAttack.hp -= damage
        }
      }
      a.age++
      if a.age > 500 && rand.Int31n(300) == 0 && a.food > 1200 {
        a.birth = true
        a.food -= 400
      } else {
        a.birth = false
      }
      if a.food == 0 {
        a.dead = true
        a.subtype.starvedCount++
        continue
      }
      if a.age > 4000 && rand.Int31n(3000) == 0 {
        a.subtype.oldAgeCount++
        a.dead = true
        continue
      }
      a.food--
    }
  }

}

func (w *world) mapModifyPhase() {
  for _, section := range w.sections {
    // Chance to grow new plant depends on number of plants

    for p := range section.plants {
      if p.dead {
        delete(section.plants, p)
      } else {
        if p.birth {
          w.addPlant(p.x-float64(rand.Int31n(50)-25), p.y-float64(rand.Int31n(50)-25))
        }
      }
    }

    for a := range section.animals {
      if a.dead {
        a.subtype.livingCount--
        delete(a.section.animals, a)
        a.section = nil
      } else {
        secPos := posInt{int32(math.Floor(a.x / float64(w.sectionsSize))), int32(math.Floor(a.y / float64(w.sectionsSize)))}
        if secPos != a.section.posInt {
          if w.sections[secPos] == nil {
            fmt.Println("Ooops!", secPos, a)
            continue
          }
          // fmt.Println("Moving from", a.section.posInt, secPos, a.id)
          delete(section.animals, a)
          w.sections[secPos].animals[a] = struct{}{}
          a.section = w.sections[secPos]
        }
        if a.birth {
          w.birth(a)
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
    for p := range section.plants {
      p.birth = false
      if rand.Int31n(80) == 0 {
        plantsCount := int32(0)
        for i := section.x - 1; i <= section.x+1; i++ {
          for j := section.y - 1; j <= section.y+1; j++ {
            s := w.sections[posInt{i, j}]
            if s == nil {
              continue
            }
            for otherPlant := range s.plants {
              if otherPlant.distance(&p.pos) < 40 {
                plantsCount++
              }
            }
          }

        }
        tropicality := p.section.y + 20
        if plantsCount-tropicality < 80 && rand.Int31n(10*plantsCount*plantsCount/(10+tropicality)+1) == 0 {
          p.birth = true
        }
      }
      if p.age > 1500 && rand.Int31n(1500) == 0 {
        p.dead = true
      } else {
        p.age++
      }

    }
  }
  wg.Done()
}

var sightRange float64 = 24

func (w *world) animalAi(a *animal) {
  // watch out, is there some food around?
  a.dMove = nil
  a.dEatPlant = nil
  a.dAttack = nil

  if a.food > 2000 {
    a.wandering = nil // just sleep
    return
  }

  var speed float64
  if a.age < 500 {
    speed = 0.75 * a.subtype.speed
  } else if a.age > 4000 {
    speed = 0.5 * a.subtype.speed
  } else {
    speed = a.subtype.speed
  }

  var closestVictim *animal = nil
  var victimDist = sightRange * 1.5
  for i := a.section.x - 1; i <= a.section.x+1; i++ {
    for j := a.section.y - 1; j <= a.section.y+1; j++ {
      s := w.sections[posInt{i, j}]
      if s != nil {
        for victim := range s.animals {
          if !a.isAlly(victim) && victim.distance(&a.pos) <= victimDist {
            victimDist = victim.distance(&a.pos)
            closestVictim = victim
          }
        }
      }
    }

  }

  if closestVictim != nil {
    a.wandering = nil
    var xChange, yChange float64
    if speed >= victimDist {
      xChange = (closestVictim.x - a.x)
      yChange = (closestVictim.y - a.y)
      a.dAttack = closestVictim
    } else {
      xChange = (closestVictim.x - a.x) * (speed / victimDist)
      yChange = (closestVictim.y - a.y) * (speed / victimDist)
    }
    a.dMove = &pos{xChange, yChange}

  } else {
    var closestPlant *plant = nil
    var plantDist = sightRange
    for i := a.section.x - 1; i <= a.section.x+1; i++ {
      for j := a.section.y - 1; j <= a.section.y+1; j++ {
        s := w.sections[posInt{i, j}]
        if s != nil {

          for p := range s.plants {
            if p.age > 100 && p.distance(&a.pos) <= plantDist {
              plantDist = p.distance(&a.pos)
              closestPlant = p
            }
          }
        }
      }

    }

    if closestPlant != nil {
      a.wandering = nil
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
      w.wander(a, speed)
    }
  }
}

func (w *world) wander(a *animal, speed float64) {
  if a.wandering == nil {
    wanderingX := a.x - float64(rand.Int31n(100)-50)
    wanderingY := a.y - float64(rand.Int31n(100)-50)
    secPos := posInt{int32(math.Floor(wanderingX / float64(w.sectionsSize))), int32(math.Floor(wanderingY / float64(w.sectionsSize)))}
    if w.sections[secPos] != nil {
      a.wandering = &pos{wanderingX, wanderingY}
    }

  } else {
    wandDist := a.wandering.distance(&a.pos)
    var xChange, yChange float64
    if speed >= wandDist {
      a.wandering = nil
    } else {
      xChange = (a.wandering.x - a.x) * (speed / wandDist)
      yChange = (a.wandering.y - a.y) * (speed / wandDist)
    }
    a.dMove = &pos{xChange, yChange}

  }
}

func (w *world) draw(pos pos, size pos) {

  for _, section := range w.sections {
    if float64((section.x+1)*w.sectionsSize) < pos.x || float64((section.x)*w.sectionsSize) > pos.x+size.x {
      continue
    }
    if float64((section.y+1)*w.sectionsSize) < pos.y || float64((section.y)*w.sectionsSize) > pos.y+size.y {
      continue
    }
    w.renderer.SetDrawColor(0, 255, 0, 255)
    for p := range section.plants {
      w.renderer.DrawPoint(int(p.x-pos.x), int(p.y-pos.y))
    }

    for a := range section.animals {
      w.renderer.SetDrawColor(a.red, a.green, a.blue, 255)
      w.renderer.DrawPoint(int(a.x-pos.x), int(a.y-pos.y))
      w.renderer.DrawPoint(int(a.x-pos.x)+1, int(a.y-pos.y))
      w.renderer.DrawPoint(int(a.x-pos.x)-1, int(a.y-pos.y))
      w.renderer.DrawPoint(int(a.x-pos.x), int(a.y-pos.y)+1)
      w.renderer.DrawPoint(int(a.x-pos.x), int(a.y-pos.y)-1)
    }
  }

}

func (p1 *pos) distance(p2 *pos) float64 {
  return math.Sqrt((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y))
}

func (w *world) createSubtypes() {
  w.subtypes = make([]*subtype, 6)
  w.subtypes[0] = &subtype{
    name:     "Red",
    speed:    1,
    red:      255,
    green:    0,
    blue:     0,
    strength: 12,
    hp:       100,
  }
  w.subtypes[1] = &subtype{
    name:     "Blue",
    speed:    1,
    red:      0,
    green:    0,
    blue:     255,
    strength: 12,
    hp:       100,
  }
  w.subtypes[2] = &subtype{
    name:     "Yellow",
    speed:    1.2,
    red:      255,
    green:    255,
    blue:     0,
    strength: 10,
    hp:       100,
  }
  w.subtypes[3] = &subtype{
    name:     "Purple",
    speed:    1.2,
    red:      255,
    green:    0,
    blue:     255,
    strength: 10,
    hp:       100,
  }
  w.subtypes[4] = &subtype{
    name:     "Teal",
    speed:    1,
    red:      0,
    green:    255,
    blue:     255,
    strength: 10,
    hp:       120,
  }
  w.subtypes[5] = &subtype{
    name:     "Silver",
    speed:    1,
    red:      220,
    green:    220,
    blue:     220,
    strength: 10,
    hp:       120,
  }
}


// Returns true if target is his friend
func (a* animal) isAlly(target *animal) bool {
  return a.subtype == target.subtype
}

func (w *world) animalAi2(a *animal) {

  // Gather inputs
  // List of animals and plants, segregated by input.
  // Resolution of the eye - maximal resolution of 60
  // Inputs are provided as a list <0,60).
  // Inputs can be nothing, friend, foe, plant
  // That gives a total of 240 inputs (pretty heavy)
  // Each of them gives also an information about distance. (? should it?)
  // There is an extra input (I'm attacked - 1 bit)
  // Also there is a memory of additional 120 inputs / outputs - connected
  // Can attack - input

  // Actions are:
  // Attack - (closest foe in a radius gets attacked). It costs some energy though.
  // Move - forward, backwards. (speed?)
  // Rotate - <0, 1) - <left, right)
  // Eat plant - (closes)

}