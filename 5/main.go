package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readInput(filename string) (map[int][]int, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot read file: %w", err)
	}

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

func main() {
	start := time.Now()
	orderMap, input, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// for k, v := range orderMap {
	// 	fmt.Printf("%d | %v\n", k, v)
	// }
	// fmt.Println()

	sum := 0
	for _, line := range input {
		if validateInput(line, orderMap) {
			middle_value := line[len(line)/2]
			sum += middle_value
		} else {
			//fmt.Println(line)
		}
	}

	duration := time.Since(start)
	fmt.Printf("Total execution time: %v\n", duration)
	fmt.Println("sum: ", sum)
}

/* Part I
Easy linear solution: T: O(n+m), M: O(n+m), 1.5ms on go with reading file
*/

/* Part II
Trying to reuse my above solution:
- Previously i had a set of numbers, that broke the rules if they appeared next,
	now if we get to such a number, we need to know why it broke the rules, so i also need a map.
- So If we get to a number that breaks a rule we need to move it before its pair, how does this affect the rules concerning that pair

- Probably if we move a number to a previous spot, we need to re-evaluate its rule.
	Not sure if its gonna be straigtforward how we update CannotAppear, after such re-order, or do we start algorithm over with clean state?
	(Not gonna worry about make it linear for now.)
*/
