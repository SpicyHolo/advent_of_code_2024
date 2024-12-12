package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Position struct {
	X, Y int
}

type PositionWithDir struct {
	pos Position
	i   int
}

func (p1 *Position) Add(p2 *Position) Position {
	return Position{p1.X + p2.X, p1.Y + p2.Y}
}

// Implement a set on a map
type Set[T comparable] map[T]struct{}

func (s Set[T]) add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) remove(element T) {
	delete(s, element)
}

func (s Set[T]) contains(element T) bool {
	_, exists := s[element]
	return exists
}

func (s Set[T]) copy() Set[T] {
	// Create a new map with the same type and initial capacity
	copy := make(Set[T], len(s))

	// Copy each key-value pair from the original map to the new map
	for key := range s {
		copy.add(key)
	}

	return copy
}

// String provides a string representation of the set.
func (s Set[T]) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	first := true
	for element := range s {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprint(element))
		first = false
	}
	builder.WriteString("}")
	return builder.String()
}

func readInput(filename string) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var res [][]byte
	for scanner.Scan() {
		res = append(res, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return res, nil
}

// Traversing helper functions
func isWithinBounds(gardenMap [][]byte, pos Position, offset int) bool {
	width, height := len(gardenMap[0])+offset, len(gardenMap)+offset
	return pos.X >= -offset && pos.X < width && pos.Y >= -offset && pos.Y < height
}

// func isWithinBounds2(gardenMap [][]byte, pos pos) bool {
// 	w, h := len(gardenMap[0]), len(gardenMap)
// 	return pos.X >= -1 &&
// 		pos.X <= w &&
// 		pos.Y >= -1 &&
// 		pos.Y <= h
// }

func getNeighbors(pos Position) []Position {
	dirs := []Position{
		{0, -1}, // UP
		{1, 0},  // RIGHT
		{0, 1},  // DOWN
		{-1, 0}, // LEFT
	}

	var neighbors []Position
	for _, dir := range dirs {
		newPos := pos.Add(&dir)
		neighbors = append(neighbors, newPos)
	}

	return neighbors
}

func getValue(gardenMap [][]byte, pos Position) byte {
	return gardenMap[pos.Y][pos.X]
}

// Use DFS to explore and mark all connnected cells within the same character region
func exploreSegment(gardenMap [][]byte, pos Position, char byte, segmentID int, segments map[Position]int) {
	// Mark position
	segments[pos] = segmentID

	// Explore neighboring positions
	for _, neighbor := range getNeighbors(pos) {
		// Skip out-ouf-bounds positions
		if !isWithinBounds(gardenMap, neighbor, 0) {
			continue
		}

		// Skip already visited positions
		if _, visited := segments[neighbor]; visited {
			continue
		}

		// If neighboring position has the same character, explore it
		if getValue(gardenMap, neighbor) == char {
			exploreSegment(gardenMap, neighbor, char, segmentID, segments)
		}
	}
}

// Find all connected segments, returns a map from all positions to segmentID
func findConnectedSegments(gardenMap [][]byte) map[Position]int {
	segments := make(map[Position]int)
	segmentID := 0
	width, height := len(gardenMap[0]), len(gardenMap)

	// Iterate over all position in the map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := Position{X: x, Y: y}

			// Skip already visited positions
			if _, visited := segments[pos]; visited {
				continue
			}

			// Start a new DFS exploration on this position
			char := getValue(gardenMap, pos)
			exploreSegment(gardenMap, pos, char, segmentID, segments)
			segmentID++
		}
	}
	return segments
}

// Given a Set of positions in a segment, find its perimeter
// By doing a dfs, and adding all neighhbors that are not in the set already
func findPerimeter(gardenMap [][]byte, segment Set[Position]) Set[PositionWithDir] {
	perimeter := make(Set[PositionWithDir])
	for pos := range segment {
		for i, neighbor := range getNeighbors(pos) {
			const offset = 1
			if !isWithinBounds(gardenMap, neighbor, offset) {
				continue
			}

			// Skip already visited positions
			if segment.contains(neighbor) {
				continue
			}

			perimeter.add(PositionWithDir{neighbor, i})
		}
	}
	return perimeter
}

func findLines(perimeter Set[PositionWithDir]) Set[PositionWithDir] {
	dirs := []Position{{1, 0}, {0, 1}}
	copy := perimeter.copy()

	for posWithDir := range perimeter {
		for _, dp := range dirs {
			next := PositionWithDir{posWithDir.pos.Add(&dp), posWithDir.i}
			if perimeter.contains(next) {
				copy.remove(posWithDir)
			}
		}
	}
	return copy
}

func main() {
	gardenMap, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find connected components
	cc := findConnectedSegments(gardenMap)

	// Get a Set containing each segment
	var segments []Set[Position]
	for pos, segmentID := range cc {
		for segmentID >= len(segments) {
			segments = append(segments, make(Set[Position]))
		}

		segments[segmentID].add(pos)
	}

	sum_part1 := 0
	sum_part2 := 0
	for _, segment := range segments {
		perimeter_set := findPerimeter(gardenMap, segment)
		line_set := findLines(perimeter_set)

		sum_part1 += len(segment) * len(perimeter_set)
		sum_part2 += len(segment) * len(line_set)
	}

	fmt.Println("Part 1: ", sum_part1)
	fmt.Println("Part 2: ", sum_part2)

}
