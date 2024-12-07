package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type operatorsMap map[byte]func(a, b int) int

// Reads input for the problem
// returns an array of results, and an array of expressions
func readInput(filename string) ([]int, [][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file: %w", err)
	}

	scanner := bufio.NewScanner(file)

	var results []int
	var expressions [][]int
	for scanner.Scan() {
		// Parse line
		line := scanner.Text()
		split := strings.Split(line, ": ")
		if len(split) != 2 {
			return nil, nil, fmt.Errorf("incorrect input, should be in format 'int: int int int ...', but got: '%v'", line)
		}
		res := split[0]
		operands := strings.Split(split[1], " ")
		if len(operands) < 2 {
			return nil, nil, fmt.Errorf("incorrect input, should be in format 'int: int int int ...', but got: '%v'", line)
		}

		// Convert to integers
		resInt, err := strconv.Atoi(res)
		if err != nil {
			return nil, nil, fmt.Errorf("incorrect input, should be in format 'int: int int int ...', but got: '%v'", line)
		}
		expr := make([]int, len(operands))
		for i, num := range operands {
			expr[i], err = strconv.Atoi(num)
			if err != nil {
				return nil, nil, fmt.Errorf("incorrect input, should be in format 'int: int int int ...', but got: '%v'", line)
			}
		}

		results = append(results, resInt)
		expressions = append(expressions, expr)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}

	return results, expressions, nil
}

// Recursively generates an array of all operator combinations of length n
func generateCombinations(byteOperatorMap operatorsMap, n int) []string {
	if n == 0 {
		return []string{""}
	}

	smallerCombinations := generateCombinations(byteOperatorMap, n-1)
	result := []string{}

	for _, combination := range smallerCombinations {
		for op := range byteOperatorMap {
			result = append(result, combination+string(op))
		}
	}
	return result
}

// Calculates the expression for every possible combination of operators.
// Returns the expression value if it matches the expected result, otherwise returns 0
func checkExpression(byteOperatorMap operatorsMap, expectedRes int, expr []int) int {
	opCombinations := generateCombinations(byteOperatorMap, len(expr)-1)

	for _, operators := range opCombinations {
		res := expr[0]
		for i, op := range operators {
			res = byteOperatorMap[byte(op)](res, expr[i+1])

			// Exit when exceeded expected result
			if res > expectedRes {
				break
			}
		}
		if res == expectedRes {
			return res
		}
	}
	return 0
}

/*
Given n operands in an expression and, m operators, given i expressions.
For each expression:
- Have to run (n-1) operations, for m^(n-1) combinations.
- So complexity is O(n*m^n), without || opperator
Part I: (2 operators, n-1 operations, i expressions), assuming less operands than expressions
O(i^2 * 2^i) -> O(n^2 * 2^n)
Memory: O(2^n), we store all the combinations

Part II: operator || takes log10(b) number, that depened on the biggest number in array
O(n^2 * 3^n * log(b)), where b is the biggest number in array
Memory: O(3^n), we store all the combination

Speed up possibilities:
Work backwards with / and - instead, dropping when number is not divisible / get a negative.
Reversing || is not so straight forward i guess.
*/
func main() {
	results, exprs, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Part I operators
	op := operatorsMap{
		'+': func(a, b int) int { return a + b },
		'*': func(a, b int) int { return a * b },
	}

	// Part II operators
	op2 := operatorsMap{
		'+': func(a, b int) int { return a + b },
		'*': func(a, b int) int { return a * b },
		'|': func(a, b int) int {
			placeValue := 1
			for b >= placeValue {
				placeValue *= 10
			}
			return a*placeValue + b
		},
	}

	sum := 0
	sum2 := 0

	// Calculates the sum of expression results, that match the expected result
	start := time.Now()
	for i, expr := range exprs {
		sum += checkExpression(op, results[i], expr)
		sum2 += checkExpression(op2, results[i], expr)
	}
	elapsed := time.Since(start)
	fmt.Printf("Took: %s\n", elapsed)
	fmt.Println("Part I")
	fmt.Println("sum: ", sum)

	fmt.Println("Part II")
	fmt.Println("sum: ", sum2)
}
