package main

import (
	"bufio"
	"day6/guard"
	"fmt"
	"os"
)

func readInput(filename string) (guard.GuardMap, int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Read map
	var resultMap guard.GuardMap
	var x, y int

	i := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row []bool
		line := scanner.Text()
		for j, char := range line {
			if char == '^' {
				x, y = j, i
			}
			row = append(row, char == '#')
		}
		resultMap = append(resultMap, row)
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("could not read file: %w", err)
	}

	return resultMap, x, y, nil
}

func main() {
	guard_map, x_init, y_init, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a guard
	g := guard.NewGuard(guard_map, x_init, y_init, "NORTH")

	visited, count := g.TracePath()
	fmt.Println("sum: ", count)

	// Reset guard
	count2 := g.CheckLoop(visited, x_init, y_init, "NORTH")
	fmt.Println("sum2: ", count2)
}
