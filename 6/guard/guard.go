package guard

import (
	"fmt"
	"strings"
)

type GuardMap [][]bool

// Deifne the direction map with (x, y) offsets
type DirectionMap map[string][2]int

var Directions DirectionMap = DirectionMap{
	"NORTH": {0, -1},
	"EAST":  {1, 0},
	"SOUTH": {0, 1},
	"WEST":  {-1, 0},
}

type Guard struct {
	Map       GuardMap
	X, Y      int
	Direction string
}

// Move in the current direction (doesnt check bounds)
func (g *Guard) MoveForward() {
	offset := Directions[g.Direction]
	g.X += offset[0]
	g.Y += offset[1]
}

// Check if there is a wall in front of the guard
func (g *Guard) IsWall() bool {
	offset := Directions[g.Direction]
	x, y := g.X+offset[0], g.Y+offset[1]
	return g.Map[y][x]
}

// Check if there is a border in front of guard
func (g *Guard) IsBorder() bool {
	max_x, max_y := len(g.Map[0]), len(g.Map)

	offset := Directions[g.Direction]
	x, y := g.X+offset[0], g.Y+offset[1]
	return x < 0 || y < 0 || x >= max_x || y >= max_y
}

// Turn 90 degrees, clockwise
func (g *Guard) Turn90CW() {
	var nextDir string

	switch g.Direction {
	case "NORTH":
		nextDir = "EAST"
	case "EAST":
		nextDir = "SOUTH"
	case "SOUTH":
		nextDir = "WEST"
	case "WEST":
		nextDir = "NORTH"
	}

	g.Direction = nextDir
}

// Constructor
func NewGuard(map_arr GuardMap, x, y int, direction string) *Guard {
	return &Guard{
		Map:       map_arr,
		X:         x,
		Y:         y,
		Direction: direction,
	}
}

// Stringer
func (g GuardMap) String() string {
	var builder strings.Builder

	builder.WriteString("Map: \n")
	for i, row := range g {
		for _, item := range row {
			if item {
				builder.WriteString("#")
			} else {
				builder.WriteString(".")
			}
		}

		if i < len(g)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// Prints guard's current state (position, heading direction, etc.)
func (g *Guard) String() string {
	var builder strings.Builder
	if g.IsBorder() {
		fmt.Fprintf(&builder, "Guard: (%v, %v) Facing: %v, ", g.X, g.Y, g.Direction)
		builder.WriteString("Finished!\n")
		return builder.String()
	}

	if g.IsWall() {
		fmt.Fprintf(&builder, "Guard: (%v, %v) Facing: %v\n", g.X, g.Y, g.Direction)
		return builder.String()
	}
	return ""
}

// Move the guard according to rules, returns false if on the edge of map (signaling end of motion)
func (g *Guard) nextMove() (notFinished bool) {
	switch {
	case g.IsBorder():
		return false

	case g.IsWall():
		g.Turn90CW()
		return true

	default:
		g.MoveForward()
		return true
	}
}

func (g *Guard) TracePath() (GuardMap, int) {
	// Initialise visited array
	w, h := len(g.Map[0]), len(g.Map)
	visited := make(GuardMap, h)
	for i := range visited {
		visited[i] = make([]bool, w)
	}

	count := 0
	for g.nextMove() {
		if !visited[g.Y][g.X] {
			count++
			visited[g.Y][g.X] = true
		}
	}

	return visited, count
}

// Part II

// Key for direction set (x, y, direction)
// Stores all unique visited locations
type posWithDir struct {
	x, y int
	dir  string
}

// Set guard's position and orientation
func (g *Guard) set_guard(x_init, y_init int, dir_init string) {
	g.X, g.Y, g.Direction = x_init, y_init, dir_init
}

// Helper function, sets one additional wall for the duration of execution, traverser the map returns True if it finds a loop.
func (g *Guard) checkLoopHelper(x, y int) bool {
	// Initialise visited array
	visited := make(map[posWithDir]struct{})

	// Set x,y to wall for the execution.
	g.Map[y][x] = true
	defer func() {
		g.Map[y][x] = false
	}()

	// Traverse the map, if we're in the same position, with the same orientation, it's a loop!
	for g.nextMove() {
		key := posWithDir{g.X, g.Y, g.Direction}
		if _, exists := visited[key]; exists {
			return true
		}

		visited[key] = struct{}{}
	}
	return false
}

// Counts the amount of possible wall locations, that create a loop.
// Searches only the previously visited positions
func (g *Guard) CheckLoop(visited GuardMap, x_init, y_init int, dir_init string) int {
	count := 0

	// For each visited position, check if adding wall will create a loop
	for y, row := range g.Map {
		for x := range row {
			if visited[y][x] {

				fmt.Println("Checking: ", x, y)
				g.set_guard(x_init, y_init, dir_init)
				if g.checkLoopHelper(x, y) {
					count++
				}
			}
		}
	}
	return count
}

// Part II:
// For each possible guard position, check for 'square' of walls, if three positions have walls, adding one will cause a loop.
/*
	- Check for wall in front of guard, than to the right of him
	- Could simulate all possible new wall placement on that line, and see if it completes a loop

	- We simulate all possible wall locations on Agent's path? adding only one wall, and seeing if it creates a loop ? (how to terminate asap on loop?)
		Turing's machine problem -> you cannot determine if program is in an infinite loop xD
*/
// Can i do this somehow without simulating guard's position?,

// Okay idea 2:
/*
	Somehow check for structures in an array that would complete a 'square', see if my agent has visited all of positions inside it, if so it's valid.
	Shouldn't be that hard to compute all posible 'squares' like that.
	Simplest way would be to first look at all possible triplets, but that is probably n!
	- Or add a square at each position and check if it creates such structures.
*/
// Idea 3:
/*
	Transfer a map from walls to all turning point locations.
	This should make computing 'loops' easy.
*/
