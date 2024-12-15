package warehouse1

import (
	. "day15/utils"
	"fmt"
	"strings"
)

type State struct {
	Map      [][]byte
	Pos      Vec2D
	Commands string
	C_ptr    int
}

func (s State) String() string {
	var builder strings.Builder
	w, h := len(s.Map[0]), len(s.Map)
	fmt.Fprintf(&builder, "Map: %dx%d, @: %v\n", w, h, s.Pos)
	if s.C_ptr < len(s.Commands) {
		fmt.Fprintf(&builder, "Next command: %c\n", s.Commands[s.C_ptr])
	}
	for _, row := range s.Map {
		builder.WriteString(string(row))
		builder.WriteByte('\n')
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
			if char != '@' {
				continue
			}

			s.Pos = *NewVec2D(x, y)
			return s, nil
		}
	}
	return nil, fmt.Errorf("no '@' representing robot, found in the map")
}

// func (s *State) isWall(pos Vec2D) bool {
// 	return ArrayGet(pos, s.Map) == '#'
// }

func (s *State) isBox(pos Vec2D) bool {
	return ArrayGet(pos, s.Map) == 'O'
}

func (s *State) isFree(pos Vec2D) bool {
	return ArrayGet(pos, s.Map) == '.'
}

func (s *State) NextCommand() bool {
	if s.C_ptr == len(s.Commands) {
		return false
	}

	// Fetch command (direction)
	command := s.Commands[s.C_ptr]
	s.C_ptr++
	dir := COMMANDS[command]

	s.move(dir)
	return true
}

func (s *State) updateMap(pos Vec2D, char byte) {
	ArraySet(pos, s.Map, char)
}

func (s *State) moveBox(pos, dir Vec2D) {
	initial_box := pos.Add(dir)
	cur_pos := initial_box

	// Moves right until founds a wall / empty spot.
	for s.isBox(cur_pos) {
		cur_pos = cur_pos.Add(dir)
	}

	// If it's empty, move entire line of boxes, else do nothing
	if s.isFree(cur_pos) {
		s.updateMap(pos, '.')
		s.updateMap(initial_box, '@')
		s.Pos = initial_box

		s.updateMap(cur_pos, 'O')
	}
}

func (s *State) move(dir Vec2D) {
	newPos := s.Pos.Add(dir)

	switch {
	case s.isFree(newPos):
		s.updateMap(s.Pos, '.')
		s.updateMap(newPos, '@')
		s.Pos = newPos
	case s.isBox(newPos):
		s.moveBox(s.Pos, dir)
	}
}

func (s *State) Score() (score int) {
	for y, row := range s.Map {
		for x, char := range row {
			if char == 'O' {
				score += y*100 + x
			}
		}
	}
	return
}
