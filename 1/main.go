package main

import (
	"bufio" // For reading file line by line
	"fmt"   // For printing
	"os"    // For opening files
	"sort"
	"strconv"
	"strings" // for string operations
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Reads an input file, returns an array of strings for each line
func getInput() ([]string, error) {
	fileName := "input.txt"

	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return data, nil
}

// Parses each line, returning two lists of ints, for left and rigth column of text.
func parseInput(data []string) (leftArray, rightArray []int, err error) {
	for _, line := range data {
		numStrs := strings.Fields(line) // Splits by whitespaces

		if len(numStrs) != 2 {
			return nil, nil, fmt.Errorf("Each line must contain exactly 2 elements, but got %d", len(numStrs))
		}

		// Convert string to integer
		left, err1 := strconv.Atoi(numStrs[0])
		right, err2 := strconv.Atoi(numStrs[1])

		if err1 != nil {
			return nil, nil, fmt.Errorf("error converting '%s'  to integer: %v", numStrs[0], err1)
		}

		if err2 != nil {
			return nil, nil, fmt.Errorf("error converting '%s'  to integer: %v", numStrs[1], err2)
		}

		// Add to respective arrays
		leftArray = append(leftArray, left)
		rightArray = append(rightArray, right)
	}

	return leftArray, rightArray, nil
}

// calculateSumOfDiffs calculates the sum of absolute differences between two arrays
func calculateSumOfDiffs(leftArray, rightArray []int) int {
	// Iterate over the arrays
	minLength := len(leftArray)
	if len(rightArray) < minLength {
		minLength = len(rightArray)
	}

	sum := 0
	for i := 0; i < minLength; i++ {
		diff := leftArray[i] - rightArray[i]
		sum += AbsInt(diff)
	}
	return sum
}

// Given an array of ints, create a map that shows frequency for each number
func createFrequencyMap(arr []int) map[int]int {
	res := make(map[int]int)

	for _, num := range arr {
		res[num]++
	}

	return res
}

func calculateSimilarityScore(arr []int, frequencyMap map[int]int) int {
	score := 0
	for _, num := range arr {
		score += num * frequencyMap[num]
	}
	return score
}

func main() {
	// Parsing input
	lines, err := getInput()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	leftArray, rightArray, err := parseInput(lines)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	/*
		PART I
		Sorting the array, should be enough to solve this.
		After sorting go over two arrays, and add their difference (don't forget the absolute value!)
		Due to sorting, solution is O(n*log(n))
	*/
	sort.Ints(leftArray)
	sort.Ints(rightArray)

	sum := calculateSumOfDiffs(leftArray, rightArray)
	fmt.Println("Sum: ", sum)

	/*
		PART II
		Create a histogram for each element in the right list (as hashmap)
		building the map takes O(m), since we need to get over the entire list
		later for each term in left list, we need to search hashmap, which is O(1) for each, total O(n)
		So solution is : O(n+m), where n, m is length of left, right list
		n = m, so solution is O(n)
	*/

	// Parse again, to avoid the arrays being sorted
	leftArray, rightArray, err = parseInput(lines)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	rightMap := createFrequencyMap(rightArray)
	score := calculateSimilarityScore(leftArray, rightMap)
	fmt.Println("Similarity Score: ", score)
}
