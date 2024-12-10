package main

import (
	"day9/file_operations"
	"day9/input"
	"fmt"
)

func part1(filename string) {
	data, err := input.ParseInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Deleted this by accident, don't want to re-implement this TODO
	//file_operations.FillEmpty()
	fmt.Println("check_sum: ", file_operations.CheckSum(data))
}

func part2(filename string) {
	data, err := input.ParseInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	file_operations.FillEmpty2(data)
	fmt.Println("check_sum: ", file_operations.CheckSum(data))
}

func main() {
	filename := "./inputs/input.txt"
	part1(filename)
	part2(filename)
}
