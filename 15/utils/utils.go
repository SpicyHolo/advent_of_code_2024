package utils

import "fmt"

const (
	UP    = 0
	RIGHT = 1
	DOWN  = 2
	LEFT  = 3
)

var DIRECTIONS = map[int]Vec2D{
	UP:    {0, -1},
	RIGHT: {1, 0},
	DOWN:  {0, 1},
	LEFT:  {-1, 0},
}

var COMMANDS = map[byte]Vec2D{
	'^': DIRECTIONS[UP],
	'>': DIRECTIONS[RIGHT],
	'v': DIRECTIONS[DOWN],
	'<': DIRECTIONS[LEFT],
}

type Vec2D struct {
	X, Y int
}

func NewVec2D(x, y int) *Vec2D {
	v := new(Vec2D)
	v.X = x
	v.Y = y
	return v
}

func (v Vec2D) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

func (v1 *Vec2D) Add(v2 Vec2D) Vec2D {
	return Vec2D{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 *Vec2D) Sub(v2 Vec2D) Vec2D {
	return Vec2D{v1.X - v2.X, v1.Y - v2.Y}
}

func (v *Vec2D) Scale(a int) Vec2D {
	return Vec2D{a * v.X, a * v.Y}
}
func ArrayGet[T any](v Vec2D, array [][]T) T {
	return array[v.Y][v.X]
}

func ArraySet[T any](v Vec2D, array [][]T, value T) {
	array[v.Y][v.X] = value
}

// Implement a set on a map
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) Remove(element T) {
	delete(s, element)
}

func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}
