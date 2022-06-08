package maps

import (
  // "errors"
  "log"
  "fmt"
  // "os"
  // "bufio"

  "github.com/BattlesnakeOfficial/rules"
)

type CoreyjaMazeMap struct{}

func init() {
  mazeMap := CoreyjaMazeMap{}
  globalRegistry.RegisterMap(mazeMap)
}

func (m CoreyjaMazeMap) ID() string {
  return "coreyja_maze"
}

func (m CoreyjaMazeMap) Meta() Metadata {
  return Metadata{
    Name:        "Coreyja Maze",
    Description: "Solo Maze where you need to find the food",
    Author:      "coreyja",
  }
}

func (m CoreyjaMazeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
  rand := settings.GetRand(initialBoardState.Turn)

  if len(initialBoardState.Snakes) != 1 {
    return rules.RulesetError("This map requires exactly one snake")
  }

  // if initialBoardState.Width != 25 || initialBoardState.Height != 25 {
  //   return rules.RulesetError("This map can only be played on a 25x25 board")
  // }

  me := initialBoardState.Snakes[0]

  tempBoardState := rules.NewBoardState(initialBoardState.Width, initialBoardState.Height)

  topRightCorner := rules.Point{X: initialBoardState.Width -1 , Y: initialBoardState.Height - 1}

  m.SubdivideRoom(tempBoardState, rand, rules.Point{X: 0, Y: 0}, topRightCorner, make([]int, 0), make([]int, 0), 0)

  editor.ClearHazards()

  for _, point := range tempBoardState.Hazards {
    editor.AddHazard(point)
  }

  editor.PlaceSnake(me.ID, []rules.Point{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}}, 100)
  editor.AddFood(topRightCorner)

  // return errors.New("We don't want to actually setup the board right now")
  return nil
}

func (m CoreyjaMazeMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
  me := lastBoardState.Snakes[0]

  if len(me.Body) >= 6  {
    // This will create a new maze
    m.SetupBoard(lastBoardState, settings, editor)

    return nil
  }

  if len(lastBoardState.Food) == 0 {
    foodPlaced := false
    tries := 0

    for (!foodPlaced) {
      rand := settings.GetRand(lastBoardState.Turn + tries)

      x := rand.Intn(lastBoardState.Width)
      y := rand.Intn(lastBoardState.Height)

      log.Print(fmt.Sprintf("Trying to place food at (%v, %v)", x, y))


      if !containsPoint(lastBoardState.Hazards, rules.Point{X: x, Y: y}) {
        editor.AddFood(rules.Point{X: x, Y: y})
        foodPlaced = true
      }

      tries++
    }
  } else {
    editor.PlaceSnake(me.ID, me.Body, 100)
  }

  return nil
}

func (m CoreyjaMazeMap) SubdivideRoom(tempBoardState *rules.BoardState, rand rules.Rand, lowPoint rules.Point, highPoint rules.Point, disAllowedHorizontal []int, disAllowedVertical []int, depth int) bool {
  didSubdivide := false

  log.Print("\n\n\n")
  log.Print(fmt.Sprintf("Subdividing room from %v to %v", lowPoint, highPoint))
  log.Print(fmt.Sprintf("disAllowedVertical %v", disAllowedVertical))
  log.Print(fmt.Sprintf("disAllowedHorizontal %v", disAllowedHorizontal))
  printMap(tempBoardState)
  fmt.Print("Press 'Enter' to continue...")
  // bufio.NewReader(os.Stdin).ReadBytes('\n')

  verticalWallPosition := -1
  horizontalWallPosition := -1
  newVerticalWall := make([]rules.Point, 0)
  newHorizontalWall := make([]rules.Point, 0)

  if (highPoint.X - lowPoint.X <= 2 && highPoint.Y - lowPoint.Y <= 2) {
    return false
  }

  // TODO: We need to make sure all the walls aren't disallowed here
  // I think thats our infinite loop problem
  verticalChoices := make([]int, 0)
  for i := lowPoint.X + 1; i < highPoint.X - 1; i++ {
    if !contains(disAllowedVertical, i) {
      verticalChoices = append(verticalChoices, i)
    }
  }
  if len(verticalChoices) > 0 {
    verticalWallPosition = verticalChoices[rand.Intn(len(verticalChoices))]
    log.Print(fmt.Sprintf("drawing Vertical Wall at %v", verticalWallPosition))

    for y := lowPoint.Y; y <= highPoint.Y; y++ {
      newVerticalWall = append(newVerticalWall, rules.Point{X: verticalWallPosition, Y: y})
    }

    didSubdivide = true
  }

  /// We can only draw a horizontal wall if there is enough space
  horizontalChoices := make([]int, 0)
  for i := lowPoint.Y + 1; i < highPoint.Y - 1; i++ {
    if !contains(disAllowedHorizontal, i) {
       horizontalChoices = append( horizontalChoices, i)
    }
  }
  if len(horizontalChoices) > 0 {
    horizontalWallPosition = horizontalChoices[rand.Intn(len(horizontalChoices))]
    log.Print(fmt.Sprintf("drawing horizontal Wall at %v", horizontalWallPosition))

    for x := lowPoint.X; x <= highPoint.X; x++ {
      newHorizontalWall = append(newHorizontalWall, rules.Point{X: x, Y: horizontalWallPosition})
    }


    didSubdivide = true
  }

  /// Here we make cuts in the walls
  if len(newVerticalWall) > 1 && len(newHorizontalWall) > 1 {
    log.Print("Need to cut with both walls")
    intersectionPoint := rules.Point{ X: verticalWallPosition, Y: horizontalWallPosition }

    newNewVerticalWall, verticalHoles := cutHoles(newVerticalWall, intersectionPoint, rand)
    newVerticalWall = newNewVerticalWall
    // disAllowedVertical = make([]int, 0)
    for _, hole := range verticalHoles {
      disAllowedHorizontal = append(disAllowedHorizontal, hole.Y)
    }
    log.Print(fmt.Sprintf("Vertical Cuts are at %v", verticalHoles))

    newNewHorizontalWall, horizontalHoles := cutHoles(newHorizontalWall, intersectionPoint, rand)
    newHorizontalWall = newNewHorizontalWall
    // disAllowedHorizontal = make([]int, 0)
    for _, hole := range horizontalHoles {
      disAllowedVertical = append(disAllowedVertical, hole.X)
    }
    log.Print(fmt.Sprintf("Horizontal Cuts are at %v", horizontalHoles))
  } else if len(newVerticalWall) > 1 {
    log.Print("Only a vertical wall needs cut")
    segmentToRemove := rand.Intn(len(newVerticalWall) - 1)
    hole := newVerticalWall[segmentToRemove]
    newVerticalWall = remove(newVerticalWall, segmentToRemove)

    disAllowedHorizontal = append(disAllowedHorizontal, hole.Y)
    log.Print(fmt.Sprintf("Cuts are at %v from index %v", hole, segmentToRemove))
  } else if len(newHorizontalWall) > 1 {
    log.Print("Only a horizontal wall needs cut")
    segmentToRemove := rand.Intn(len(newHorizontalWall) - 1)
    hole := newHorizontalWall[segmentToRemove]
    newHorizontalWall = remove(newHorizontalWall, segmentToRemove)

    disAllowedVertical = append(disAllowedVertical, hole.X)
    log.Print(fmt.Sprintf("Cuts are at %v from index %v", hole, segmentToRemove))
  }

  for _, point := range newVerticalWall {
    tempBoardState.Hazards = append(tempBoardState.Hazards, point)
  }
  for _, point := range newHorizontalWall {
    tempBoardState.Hazards = append(tempBoardState.Hazards, point)
  }

  /// We have both so need 4 sub-rooms
  if (verticalWallPosition != -1 && horizontalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth + 1)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth + 1)

    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: verticalWallPosition, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth + 1)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth + 1)
  } else if (verticalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth + 1)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth + 1)
  } else if (horizontalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth + 1)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth + 1)
  }

  return didSubdivide
}

func remove(s []rules.Point, i int) []rules.Point {
  s[i] = s[len(s)-1]
  return s[:len(s)-1]
}

func pos(s []rules.Point, e rules.Point) int {
  for i, a := range s {
    if a == e {
      return i
    }
  }
  return -1
}

/// Return value is first the wall that has been cut, the second is the holes we cut out
func cutHoles(s []rules.Point, intersection rules.Point, rand rules.Rand) ([]rules.Point, []rules.Point) {
  holes := make([]rules.Point, 0)

  index := pos(s, intersection)

  if index != 0 {
    firstSegmentToRemove := rand.Intn(index)
    holes = append(holes, s[firstSegmentToRemove])
    s  = remove(s, firstSegmentToRemove)
  }

  index = pos(s, intersection)
  if index != len(s) - 1 {
    secondSegmentToRemove := rand.Intn(len(s) - index - 1) + index + 1

    holes = append(holes, s[secondSegmentToRemove])
    s = remove(s, secondSegmentToRemove)
  }

  return s, holes
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func containsPoint(s []rules.Point, e rules.Point) bool {
    for _, a := range s {
        if a.X == e.X && a.Y == e.Y {
            return true
        }
    }
    return false
}
