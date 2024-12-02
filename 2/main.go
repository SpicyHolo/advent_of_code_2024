package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func parseInput(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var res [][]int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		numStrs := strings.Fields(scanner.Text())
		nums := make([]int, len(numStrs))

		for i, numStr := range numStrs {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("error converting '%s' to integer: %v", numStr, err)
			}
			nums[i] = num
		}

		res = append(res, nums)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Could not read file: %w", err)
	}

	return res, nil
}

func hasChangedMonotonicity(diff int, is_decreasing bool) bool {
	return (diff > 0 && !is_decreasing) || (diff < 0 && is_decreasing)
}

// Count the amount of safe rows
func countSafe(data [][]int) int {
	counter := 0

	for _, row := range data {
		if isSafe(row, false) {
			counter++
		}
	}
	return counter
}

// Calculate safe row, by removing each element and checking whether result is safe
func isSafeRemove(row []int) bool {
	for i := range row {
		newRow := removeElement(row, i)
		if isSafe(newRow, true) {
			return true
		}
	}
	return false
}

// Given a row of integers, check if its safe
func isSafe(row []int, alreadyRemoved bool) bool {

	// Determine initial monotonicity
	first_diff := row[0] - row[1]
	if AbsInt(first_diff) < 1 || AbsInt(first_diff) > 3 {
		if alreadyRemoved {
			return false
		} else {
			return isSafeRemove(row)
		}
	}
	is_decreasing := first_diff > 0

	// Loop over the row
	for i := 1; i < len(row)-1; i++ {
		diff := row[i] - row[i+1]

		// Conditions for safe data row, difference should be in (1, 3), and series should be monotonous
		if AbsInt(diff) < 1 || AbsInt(diff) > 3 || hasChangedMonotonicity(diff, is_decreasing) {
			if alreadyRemoved {
				return false
			} else {
				return isSafeRemove(row)
			}
		}
	}
	return true
}

// Removes element at index from array, creates a copy.
func removeElement(arr []int, index int) []int {
	// Create a copy of the slice before removing the element
	newArr := append([]int(nil), arr...) // Create a new slice with the same elements
	// Now remove the element at the specified index from the copied slice
	return append(newArr[:index], newArr[index+1:]...)
}

func main() {
	// Parse the input file
	// data, err := parseInput("test_input.txt")
	data, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Count and print the number of safe rows
	count := countSafe(data)
	fmt.Println("Count: ", count)
}
