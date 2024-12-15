package main

import (
	"bufio"
	wh1 "day15/warehouse1"
	wh2 "day15/warehouse2"
	"fmt"
	"os"
	"strings"
)

func parseInput(filename string) (*wh1.State, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	data := string(buf)

	split := strings.Split(data, "\r\n\r\n")
	warehouse := strings.Split(split[0], "\r\n")
	commands := strings.ReplaceAll(split[1], "\r\n", "")

	// Convert warehouse to [][]byte
	warehouse_bytes := make([][]byte, len(warehouse))
	for i, s := range warehouse {
		warehouse_bytes[i] = []byte(s)
	}

	state, err := wh1.NewState(warehouse_bytes, commands)
	if err != nil {
		return nil, fmt.Errorf("error loading state: %w", err)
	}

	return state, nil
}

func parseInput2(state *wh1.State) (*wh2.State, error) {
	warehouse, commands := state.Map, state.Commands

	newMap := make([][]byte, len(warehouse))

	for i, row := range warehouse {
		temp := string(row)
		temp = strings.ReplaceAll(temp, "#", "##")
		temp = strings.ReplaceAll(temp, "O", "[]")
		temp = strings.ReplaceAll(temp, ".", "..")
		temp = strings.ReplaceAll(temp, "@", "@.")
		newMap[i] = []byte(temp)
	}

	new_state, err := wh2.NewState(newMap, commands)
	if err != nil {
		return nil, err
	}
	return new_state, nil
}

func example() {
	state, err := parseInput("test_input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*state)
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	for state.NextCommand() {
		fmt.Println(*state)
	}
	fmt.Println("Score: ", state.Score())
}

func part1() {
	state, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*state)

	for state.NextCommand() {
		//fmt.Println(*state)
	}
	fmt.Println("Score: ", state.Score())
}

func example2(slow bool) {
	state, err := parseInput("test_input.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	newState, err := parseInput2(state)
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(*newState)
	if slow {
		scanner.Scan()
	}

	for newState.NextCommand() {
		if slow {
			fmt.Println(*newState)
			scanner.Scan()
		}
	}

	fmt.Println(*newState)
	fmt.Println("Score: ", newState.Score())
}

func part2(slow bool) {
	state, err := parseInput("input.txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	newState, err := parseInput2(state)
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(*newState)
	if slow {
		scanner.Scan()
	}

	for newState.NextCommand() {
		if slow {
			fmt.Println(*newState)
			scanner.Scan()
		}
	}

	fmt.Println(*newState)
	fmt.Println("Score: ", newState.Score())
}

func main() {
	// part1()

	example2(false)
	part2(false)
}
