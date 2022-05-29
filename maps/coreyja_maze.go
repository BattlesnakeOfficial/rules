package maps

import (
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
  rand := settings.GetRand(0)

  if len(initialBoardState.Snakes) != 1 {
    return rules.RulesetError("This map requires exactly one snake")
  }

  if initialBoardState.Width != 25 || initialBoardState.Height != 25 {
    return rules.RulesetError("This map can only be played on a 25x25 board")
  }

  me := initialBoardState.Snakes[0]
  editor.PlaceSnake(me.ID, []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}}, 100)

  tempBoardState := rules.NewBoardState(initialBoardState.Width, initialBoardState.Height)

  m.SubdivideRoom(tempBoardState, rand, rules.Point{X: 0, Y: 0}, rules.Point{X: 24, Y: 24}, make([]int, 0), make([]int, 0))

  editor.ClearHazards()

  for _, point := range tempBoardState.Hazards {
    editor.AddHazard(point)
  }

  editor.AddFood(rules.Point{X: 24, Y: 24})

  return nil
}

func (m CoreyjaMazeMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
  // me := lastBoardState.Snakes[0]

  // if len(lastBoardState.Food) == 0 {
  //   editor.PlaceSnake(me.ID, []rules.Point{{X: 0, Y: 1}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}}, 100)

  //   if len(lastBoardState.Snakes[0].Body) % 2 == 0 {
  //     editor.AddFood(rules.Point{X: 20, Y: 20})
  //   } else {
  //     editor.AddFood(rules.Point{X: 3, Y: 3})
  //   }
  // }

  return nil
}

func (m CoreyjaMazeMap) SubdivideRoom(tempBoardState *rules.BoardState, rand rules.Rand, lowPoint rules.Point, highPoint rules.Point, disAllowedHorizontal []int, disAllowedVertical []int) bool {
  didSubdivide := false
  /// We can only draw a vertical wall if there is enough space

  verticalWallPosition := -1
  horizontalWallPosition := -1
  newVerticalWall := make([]rules.Point, 0)
  newHorizontalWall := make([]rules.Point, 0)

  if (highPoint.X - lowPoint.X >= 4) {
    for (verticalWallPosition == -1 || contains(disAllowedVertical, verticalWallPosition)) {
      verticalWallPosition = rand.Intn(highPoint.X - lowPoint.X - 2) + lowPoint.X + 1
    }

    for y := lowPoint.Y; y <= highPoint.Y; y++ {
      newVerticalWall = append(newVerticalWall, rules.Point{X: verticalWallPosition, Y: y})
    }

    didSubdivide = true
  }

  /// We can only draw a horizontal wall if there is enough space
  if (highPoint.Y - lowPoint.Y >= 4) {
    for (horizontalWallPosition == -1 || contains(disAllowedHorizontal, horizontalWallPosition)) {
      horizontalWallPosition = rand.Intn(highPoint.Y - lowPoint.Y - 2) + lowPoint.Y + 1
    }

    for x := lowPoint.X; x <= highPoint.X; x++ {
      newHorizontalWall = append(newHorizontalWall, rules.Point{X: x, Y: horizontalWallPosition})
    }


    didSubdivide = true
  }

  /// Here we make cuts in the walls
  if len(newVerticalWall) > 1 && len(newHorizontalWall) > 1 {
    intersectionPoint := rules.Point{ X: verticalWallPosition, Y: horizontalWallPosition }

    newNewVerticalWall, verticalHoles := cutHoles(newVerticalWall, intersectionPoint, rand)
    newVerticalWall = newNewVerticalWall
    disAllowedVertical = make([]int, 0)
    for _, hole := range verticalHoles {
      disAllowedVertical = append(disAllowedVertical, hole.X)
    }

    newNewHorizontalWall, horizontalHoles := cutHoles(newHorizontalWall, intersectionPoint, rand)
    newHorizontalWall = newNewHorizontalWall
    disAllowedHorizontal = make([]int, 0)
    for _, hole := range horizontalHoles {
      disAllowedHorizontal = append(disAllowedHorizontal, hole.Y)
    }
  } else if len(newVerticalWall) > 1 {
    segmentToRemove := rand.Intn(len(newVerticalWall))
    newVerticalWall = remove(newVerticalWall, segmentToRemove)
  } else if len(newHorizontalWall) > 1 {
    segmentToRemove := rand.Intn(len(newHorizontalWall))
    newHorizontalWall = remove(newHorizontalWall, segmentToRemove)
  }


  for _, point := range newVerticalWall {
    tempBoardState.Hazards = append(tempBoardState.Hazards, point)
  }
  for _, point := range newHorizontalWall {
    tempBoardState.Hazards = append(tempBoardState.Hazards, point)
  }

  /// We have both so need 4 sub-rooms
  if (verticalWallPosition != -1 && horizontalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition - 1}, disAllowedHorizontal, disAllowedVertical)

    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: verticalWallPosition - 1, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical)
  } else if (verticalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical)
  } else if (horizontalWallPosition != -1) {
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical)
    m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical)
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

    firstSegmentToRemove := rand.Intn(index)
    holes = append(holes, s[firstSegmentToRemove])
    s  = remove(s, firstSegmentToRemove)

    if index != 0 && index != len(s) - 1 {
      index = pos(s, intersection)
      secondSegmentToRemove := rand.Intn(len(s) - index ) + index

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
