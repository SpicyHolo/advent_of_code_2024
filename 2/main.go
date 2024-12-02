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

func countSafe(data [][]int) int {
	counter := 0

	for _, row := range data {
		is_decreasing := false
		is_safe := true

		// Determine initial monotonicity
		first_diff := row[0] - row[1]
		if AbsInt(first_diff) < 1 || AbsInt(first_diff) > 3 {
			continue
		}
		is_decreasing = first_diff > 0

		// Loop over each integer
		for i := 1; i < len(row)-1; i++ {
			diff := row[i] - row[i+1]
			diffAbs := AbsInt(diff)

			// Conditions for safe data row, difference should be in {1, 2, 3}, and series should be monotonous
			if diffAbs < 1 || diffAbs > 3 || hasChangedMonotonicity(diff, is_decreasing) {
				is_safe = false
				break
			}
		}

		if is_safe {
			counter++
		}
	}
	return counter
}

func main() {
	// Parse the input file
	data, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Count and print the number of safe rows
	count := countSafe(data)
	fmt.Println("Count: ", count)
}
