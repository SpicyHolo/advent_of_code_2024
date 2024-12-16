package warehouse1

import (
	"day15/utils"
	"fmt"
	"strings"
)

type State struct {
	Map      [][]byte
	Pos      utils.Vec2D
	Commands string
	C_ptr    int
}

func (s State) String() string {
	var builder strings.Builder
	w, h := len(s.Map[0]), len(s.Map)
	fmt.Fprintf(&builder, "Map: %dx%d, @: %v\n", w, h, s.Pos)
	if s.C_ptr < len(s.Commands) {
		fmt.Fprintf(&builder, "Next command: %c, (%d/%d)\n", s.Commands[s.C_ptr], s.C_ptr+1, len(s.Commands))
	}
	for i, row := range s.Map {
		builder.WriteString(string(row))
		if i < len(s.Map)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func NewState(Map [][]byte, commands string) (*State, error) {
	s := new(State)
	s.Map = Map
	s.Commands = commands
	s.C_ptr = 0

	// Find starting position
	for y, row := range Map {
		for x, char := range row {
			switch char {
			case '@':
				s.Pos = utils.Vec2D{X: x, Y: y}
			}
		}
	}

	return s, nil
}

// Create the copy of the entire map
func deepCopy(original [][]byte) [][]byte {
	// Create a new [][]byte slice with the same length as the original
	copied := make([][]byte, len(original))

	// Loop through each row of the original slice
	for i := range original {
		// Create a new slice for each row and copy its contents
		copied[i] = append([]byte(nil), original[i]...)
	}

	return copied
}

// Checkers
func (s *State) isBox(pos utils.Vec2D) bool {
	char := s.getValue(pos)
	return char == '[' || char == ']'
}

func (s *State) isFree(pos utils.Vec2D) bool {
	return s.getValue(pos) == '.'
}

func (s *State) getValue(pos utils.Vec2D) byte {
	return s.Map[pos.Y][pos.X]
}

func (s *State) setValue(pos utils.Vec2D, char byte) {
	s.Map[pos.Y][pos.X] = char
}

// Update's robot position (also visually)
func (s *State) updateRobot(dir utils.Vec2D) {
	newPos := s.Pos.Add(dir)
	s.setValue(s.Pos, '.')
	s.setValue(newPos, '@')
	s.Pos = newPos
}

// Updates boxes and robot position after move
/*
First we keep an array for everything to move,
For each we add the stuff behind it in movement's direction
	- If it's a box we add it's second part
	- If it's a wall, it means the move is impossible, so return false.
After that we move everything that's in the created array.const
Works for both horizontal and vertical movement <3
*/
func (s *State) updateEverything(dir utils.Vec2D) bool {
	toUpdate := []utils.Vec2D{s.Pos}
	visited := make(utils.Set[utils.Vec2D]) // Keep a set to only visit each location once, faster than searching through toUpdate array
	i := 0

	// Go over the FIFO queue
	for i < len(toUpdate) {
		// Fetch, calculate next position in movement's directoin
		cur_pos := toUpdate[i]
		new_pos := cur_pos.Add(dir)

		// Case: box
		if s.isBox(new_pos) {
			// Add box if not in queue
			if !visited.Contains(new_pos) {
				toUpdate = append(toUpdate, new_pos)
				visited.Add(new_pos)
			}

			// Add box's second part
			if s.getValue(new_pos) == '[' {
				if right := new_pos.Add(utils.DIRECTIONS[utils.RIGHT]); !visited.Contains(right) {
					toUpdate = append(toUpdate, right)
					visited.Add(right)
				}
			}
			if s.getValue(new_pos) == ']' {
				if left := new_pos.Add(utils.DIRECTIONS[utils.LEFT]); !visited.Contains(left) {
					toUpdate = append(toUpdate, left)
					visited.Add(left)
				}
			}

			// Case: wall, impossible move, so return false
		} else if s.getValue(new_pos) == '#' {
			return false
		}
		// Next instruction
		i++
	}

	// Move everything, makre sure to make a copy to avoid overwriting stuff.
	map_copy := deepCopy(s.Map)

	for _, pos := range toUpdate {
		map_copy[pos.Y][pos.X] = '.'
	}
	for _, pos := range toUpdate {

		newPos := pos.Add(dir)
		if s.Map[pos.Y][pos.X] == '@' {
			s.Pos = newPos
		}
		map_copy[newPos.Y][newPos.X] = s.Map[pos.Y][pos.X]

	}

	// Swap to new map
	s.Map = map_copy
	return true
}

// Moves the robot, and boxes
func (s *State) move(dir utils.Vec2D) {
	newPos := s.Pos.Add(dir)
	switch {

	// If posiiton is free, we just move the robot
	case s.isFree(newPos):
		s.updateRobot(dir)

	// If position is a box, we try to move it, and move the robot.
	case s.isBox(newPos):
		s.updateEverything(dir)
	}
}

// Public methods

// Takes next commandi and executes it, returns false if no command left
func (s *State) NextCommand() bool {
	if s.C_ptr == len(s.Commands) {
		return false
	}

	// Fetch command (direction)
	command := s.Commands[s.C_ptr]
	s.C_ptr++
	dir := utils.COMMANDS[command]

	s.move(dir)
	return true
}

// Calculate score according to rules
func (s *State) Score() (score int) {
	for y, row := range s.Map {
		for x, char := range row {
			if char == '[' {
				score += y*100 + x
			}
		}
	}
	return
}
