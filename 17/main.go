package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Computer struct {
	Reg  map[string]int
	Prog []int
	Ptr  int
}

func (c Computer) String() string {
	var builder strings.Builder
	builder.WriteString("Registers ")
	fmt.Fprintf(&builder, "Registers | A: %v | B: %v | C: %v\n", c.Reg["A"], c.Reg["B"], c.Reg["C"])

	if c.Ptr < len(c.Prog) {
		fmt.Fprintf(&builder, "Program: %d, next-> %v\n", c.Ptr, OPCODE2OPERATION[c.Prog[c.Ptr]])
	}

	fmt.Fprintf(&builder, "%v, len(%v)\n", c.Prog, len(c.Prog))

	return builder.String()
}

var OPCODE2OPERATION = map[int]string{
	0: "adv",
	1: "bxl",
	2: "bst",
	3: "jnz",
	4: "bxc",
	5: "out",
	6: "bdv",
	7: "cdv",
}

func parseInput(filename string) (map[string]int, []int, error) {
	registers := make(map[string]int)

	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read file: %w", err)

	}

	regex_reg := regexp.MustCompile(`Register [A-C]: (\d+)`)
	regex_program := regexp.MustCompile(`Program: ((?:\d,|\d$)+)`)

	data := string(buf)
	matches_reg := regex_reg.FindAllStringSubmatch(data, -1)
	matches_program := regex_program.FindAllStringSubmatch(data, -1)

	a, err1 := strconv.Atoi(matches_reg[0][1])
	if err1 != nil {
		return nil, nil, fmt.Errorf("could not convert to int: %w", err1)
	}
	b, err2 := strconv.Atoi(matches_reg[1][1])
	if err2 != nil {
		return nil, nil, fmt.Errorf("could not convert to int: %w", err2)
	}
	c, err3 := strconv.Atoi(matches_reg[2][1])
	if err3 != nil {
		return nil, nil, fmt.Errorf("could not convert to int: %w", err3)
	}
	registers["A"] = a
	registers["B"] = b
	registers["C"] = c

	program_nums := strings.Split(matches_program[0][1], ",")
	program := make([]int, len(program_nums))
	for i, num := range program_nums {
		numInt, err := strconv.Atoi(num)
		if err != nil {
			return nil, nil, fmt.Errorf("could not convert to int: %w", err)
		}
		program[i] = numInt
	}

	return registers, program, nil
}

func (c *Computer) combo(operand int) (res int) {
	switch operand {
	case 4:
		res = c.Reg["A"]
	case 5:
		res = c.Reg["B"]
	case 6:
		res = c.Reg["C"]
	default:
		res = operand
	}
	return
}

func (c *Computer) execute(op, operand int) (int, bool) {
	switch op {
	case 0: //adv
		c.Reg["A"] = c.Reg["A"] >> c.combo(operand)
	case 1: //bxl not tested
		c.Reg["B"] = c.Reg["B"] ^ operand
	case 2: //bst
		c.Reg["B"] = c.combo(operand) % 8
	case 3: //jnz
		if c.Reg["A"] != 0 {
			c.Ptr = operand // -2 cancels the +=2 that will be done.
			return -1, true
		}
	case 4: //bxc
		c.Reg["B"] = c.Reg["B"] ^ c.Reg["C"]
	case 5: //out
		return c.combo(operand) % 8, false
	case 6: //bdv
		c.Reg["B"] = c.Reg["A"] >> c.combo(operand)
	case 7: //cdv
		c.Reg["C"] = c.Reg["A"] >> c.combo(operand)
	}

	return -1, false
}

func (c *Computer) nextCommand() int {
	if c.Ptr >= len(c.Prog)-1 { // Halt
		return -2
	}

	op := c.Prog[c.Ptr]
	operand := c.Prog[c.Ptr+1]
	res, skipIncPtr := c.execute(op, operand)

	if !skipIncPtr {
		c.Ptr += 2
	}

	return res
}

func outputToStr(output []int) string {
	outputStr := make([]string, len(output))
	for i, num := range output {
		outputStr[i] = strconv.Itoa(num)
	}

	return strings.Join(outputStr, ",")
}

func solve(filename string) {
	// Parse Input
	reg, program, err := parseInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialise
	computer := Computer{reg, program, 0}
	var output []int

	// Modify register A
	if len(os.Args) > 2 {
		newA, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("error converting arg 2 to string: %w", err)
		}
		computer.Reg["A"] = newA
	}

	// Execute program
	for {
		fmt.Println(computer)
		res := computer.nextCommand()
		// No commands left
		if res == -2 {
			break
		}

		// Output detected
		if res != -1 {
			output = append(output, res)
		}
	}

	if slices.Equal(output, computer.Prog) {
		fmt.Println("output == program")
	}

	fmt.Println("Output: ", outputToStr(output))
}

func solve2(filename string) {
	// Parse Input
	reg, program, err := parseInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialise
	computer := Computer{reg, program, 0}
	var output []int

	// Modify register A
	if len(os.Args) > 2 {
		newA, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("error converting arg 2 to string: %w", err)
		}
		computer.Reg["A"] = newA
	}

	output = solver(computer.Reg["A"], computer.Reg["B"], computer.Reg["C"])

	if slices.Equal(output, computer.Prog) {
		fmt.Println("output == program")
	}

	fmt.Println("Program: \t", outputToStr(program))
	fmt.Println("Output: \t", outputToStr(output))
}
func solver(A, B, C int) (sol []int) {
	// Parse Input
	// reg, program, err := parseInput("input.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	for A > 0 {
		B = B % 8
		B = B ^ 1
		C = A >> B
		B = B ^ 5
		A = A >> 3
		B = B ^ C
		sol = append(sol, B%8)
	}
	print(B)
	return
}

func part2(B, C int, program []int) int {
	// A := 0
	for i, code := range program {
		fmt.Println("i: ", i)
		B = B ^ 1
		c_shift := B

		B = B ^ 5
		C = B ^ code
		temp_a := C << c_shift
		fmt.Println(temp_a)
	}
	return 0
}

// }
func main() {
	if len(os.Args) > 1 {
		solve2(os.Args[1])
		return
	}
	fmt.Println("not enough arguments, expected filename.")
}
