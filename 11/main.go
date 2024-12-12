package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func parseInput(input string) ([]int, error) {
	split := strings.Split(input, " ")
	numbers := make([]int, len(split))
	for i, num := range split {
		numInt, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("error converting to integer: %w", err)
		}
		numbers[i] = numInt
	}
	return numbers, nil
}

// Given an integer (with even no. of digits), splits into two numbers by slicing its digits
func splitStone(stone int) (int, int, error) {
	str := strconv.Itoa(stone)
	middle := len(str) / 2
	str1, str2 := str[:middle], str[middle:]
	stone1, err1 := strconv.Atoi(str1)
	stone2, err2 := strconv.Atoi(str2)
	if err1 != nil {
		return 0, 0, fmt.Errorf("error converting to integer: %w", err1)
	}
	if err2 != nil {
		return 0, 0, fmt.Errorf("error converting to integer: %w", err2)
	}

	return stone1, stone2, nil
}

// Return number of digits in an integer
func numberOfDigits(stone int) int {
	if stone == 0 {
		return 1
	}
	if stone < 0 {
		stone = -stone
	}
	return int(math.Floor(math.Log10(float64(stone)))) + 1
}

// Solve using Dynamic programming
// Array of maps, each representing time step.
// For each unique stone, stores its number.
// To update to next time step, just apply rules to all unique stones in its previous array, and pass its count to the next iteration.
func blinkNTimes(initStones []int, n_steps int) (int, error) {
	dp := make([]map[int]int, n_steps+1)
	for i := 0; i <= n_steps; i++ {
		dp[i] = make(map[int]int)
	}

	for _, stone := range initStones {
		dp[0][stone]++
	}

	for i := 0; i < n_steps; i++ {
		for stone, count := range dp[i] {
			if stone == 0 {
				dp[i+1][1] += count
			} else if numberOfDigits(stone)%2 == 0 {
				left, right, err := splitStone(stone)
				if err != nil {
					return 0, err
				}
				dp[i+1][left] += count
				dp[i+1][right] += count
			} else {
				dp[i+1][stone*2024] += count
			}
		}
	}
	totalStones := 0
	for _, count := range dp[n_steps] {
		totalStones += count
	}

	return totalStones, nil
}

func main() {
	input := "92 0 286041 8034 34394 795 8 2051489"
	stones, err := parseInput(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	count1, err := blinkNTimes(stones, 25)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Part I: ", count1)

	count2, err := blinkNTimes(stones, 75)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Part II: ", count2)
}
