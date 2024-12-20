package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func parseInput(filename string) ([]string, []string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file: %w", err)
	}

	blocks := strings.Split(string(fileContent), "\n\n")
	if len(blocks) != 2 {
		fmt.Println(len(blocks), blocks)
		return nil, nil, fmt.Errorf("invalid input file format, should contain patterns, designs seperated by an empty line.")
	}

	re := regexp.MustCompile(`[a-z]+`)

	patterns := re.FindAllString(blocks[0], -1)
	if patterns == nil {
		return nil, nil, fmt.Errorf("Invalid pattern format, should contain patterns [a-z]+ seperated by a comma.")
	}

	designs := re.FindAllString(blocks[1], -1)
	if designs == nil {
		return nil, nil, fmt.Errorf("Invalid design format, should contain design [a-z] strings, each in a new line.")
	}

	return patterns, designs, nil
}

func possible(patterns []string, design string) bool {
	if design == "" {
		return true
	}

	// Check each pattern
	for _, pattern := range patterns {
		if len(design) >= len(pattern) && design[:len(pattern)] == pattern {
			if possible(patterns, design[len(pattern):]) {
				return true
			}
		}
	}

	return false
}

// Given a list of substrings, count all the possible ways of arraning them to form the given string
func countCombinations(patterns []string, design string, cache map[string]int) int {
	// Edge case, if string is empty, only one way to arrange it.
	if design == "" {
		return 1
	}

	// Check if we already calculted the value for that string, in cache
	if val, exists := cache[design]; exists {
		return val
	}

	count := 0
	// For each substring, subtract it from the given string, and recursively call on new string (with prefix removed)
	for _, pattern := range patterns {
		if len(design) >= len(pattern) && design[:len(pattern)] == pattern {
			count += countCombinations(patterns, design[len(pattern):], cache)

		}
	}

	// Add result to cache
	cache[design] = count
	return count
}

func part1() {
	patterns, designs, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, design := range designs {
		if possible(patterns, design) {
			sum++
		}
	}
	fmt.Println(sum)
}

func part2() {
	patterns, designs, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, design := range designs {
		cache := make(map[string]int)
		sum += countCombinations(patterns, design, cache)
	}
	fmt.Println(sum)
}

func main() {
	part2()
}
