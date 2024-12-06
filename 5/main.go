package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(filename string) (map[int][]int, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot read file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Create a map
	orderMap := make(map[int][]int)
	for scanner.Scan() {
		line := scanner.Text()
		// Empty lines represents the end of defining rules, moving to input
		if line == "" {
			break
		}

		tokens := strings.Split(line, "|")
		if len(tokens) != 2 {
			return nil, nil, fmt.Errorf("Error reading the input, rule should be in a format 'int|int', but got: %v", line)
		}
		X, err1 := strconv.Atoi(tokens[0])
		Y, err2 := strconv.Atoi(tokens[1])
		if err1 != nil {
			return nil, nil, fmt.Errorf("Error reading the input, rule should be in a format 'int|int', but got: %v, %w", line, err1)
		}
		if err2 != nil {
			return nil, nil, fmt.Errorf("Error reading the input, rule should be in a format 'int|int', but got: %v, %w", line, err2)
		}

		orderMap[Y] = append(orderMap[Y], X)
	}

	// Read input
	var input [][]int
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			return nil, nil, fmt.Errorf("Found an empty line in the input part, input should be a continous paragraph.")
		}
		tokens := strings.Split(line, ",")

		nums := make([]int, len(tokens))
		for i, token := range tokens {
			num, err := strconv.Atoi(token)
			if err != nil {
				return nil, nil, fmt.Errorf("Error converting one of the tokens to int: %w", err)
			}
			nums[i] = num
		}
		input = append(input, nums)
	}

	return orderMap, input, nil
}

// m - num of rules
// n - num of input
/*
	- Reading file, creating map, input: O(n+m)

	Going through each of nums O(n)
	Searching cannotAppear, appeared set, worst-case O(n), average O(1)
	Searching orderMap set, worst case O(m), average O(1)

	Average: O(n+m), worst-case either O(n*m) or O(n^2) depending whether there are more rules or nums in input

	Algorithm:
	pre: Create a map for all map[Y] X

	1. Go through each number in an update
	2. Add the number to Appeared set
	3. Check if it's in CannotAppear set, if it is result is False
	4. Check in rules map, if so add every number that had to appear before it, and didn't to CannotAppear set
		(Those numbers cannot appear, because according to rules they had to appear before current number, but didn't)
	5. If gone through the whole list, return True

*/
func validateInput(input []int, orderMap map[int][]int) bool {
	cannotAppear := make(map[int]bool)
	appeared := make(map[int]bool)

	for _, num := range input {
		appeared[num] = true

		// Check if this number can appear
		if cannotAppear[num] {
			return false
		}

		// Check if any numbers need to appear before it,
		pairs, ok := orderMap[num]
		if ok {
			// If pair needs to appear before num, check if it did, if not add it to cannotAppear
			for _, pair := range pairs {
				if !appeared[pair] {
					cannotAppear[pair] = true
				}
			}

		}
	}
	return true
}

// Part II solution
// O(n(n+m)), assuming less rules than 'updates' -> O(n^2)
/*
	Brute Force, (I'm tired boss...)
	Use the same algorithm as in Part I, but
	but now under each key in cannotAppear we save and index of the number that inforced that rule.
	If we reach such number, we swap the rule breaking number, with rule enforcing number, and
	call the function on the new array.
	Do this, until the entire array is valid, and return the result array.

	Probably can try a linear solution, where you need to decide how to update catnnotAppear, appeared maps, after the swap.
	Redoing the rule checking for swapped number (since after the swap, more numbers come after it), could do the trick. (Will re-visit TODO)
*/
func partTwo(input []int, orderMap map[int][]int) []int {
	cannotAppear := make(map[int]int)
	appeared := make(map[int]bool)

	for i, num := range input {
		appeared[num] = true

		// Check if this number can appear
		_, ok := cannotAppear[num]
		if ok {
			// Swap and call PartTwo with new array
			j := cannotAppear[num]
			input[i], input[j] = input[j], input[i]
			return partTwo(input, orderMap)
		}

		// Check if any numbers need to appear before it,
		pairs, ok := orderMap[num]
		if ok {
			// If pair needs to appear before num, check if it did, if not add it to cannotAppear
			for _, pair := range pairs {
				if !appeared[pair] {
					_, ok := cannotAppear[num]
					if !ok {
						cannotAppear[pair] = i
					}
				}
			}

		}
	}
	return input
}

func main() {
	// Read inputs, create map of rules
	orderMap, input, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Rules:")
	for v, k := range orderMap {
		fmt.Printf("%d | %v\n", k, v)
	}

	fmt.Println("Num of. Values: ", len(input))
	fmt.Println()

	// Incorrect inputs for the second part
	var incorrect_inputs [][]int

	// Part I solution, O(n + m), assuming less rules than 'updates' -> O(n)
	sum := 0
	for _, line := range input {
		if validateInput(line, orderMap) {
			middle_value := line[len(line)/2]
			sum += middle_value
		} else {
			incorrect_inputs = append(incorrect_inputs, line)
		}
	}

	fmt.Println("Part I sum: ", sum)

	sum = 0
	for _, line := range incorrect_inputs {
		fixed_line := partTwo(line, orderMap)
		middle_value := fixed_line[len(fixed_line)/2]
		sum += middle_value
	}
	fmt.Println("Part II sum:", sum)
}
