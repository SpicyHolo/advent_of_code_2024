package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Position struct {
	X, Y int
}

func (p Position) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "(%v, %v)", p.X, p.Y)
	return builder.String()
}

func (p1 Position) Add(p2 Position) Position {
	return Position{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

type Game struct {
	ButtonA, ButtonB, PrizePos Position
}

func (g Game) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Game {A: %v, B: %v, Prize: %v}", g.ButtonA, g.ButtonB, g.PrizePos)
	return builder.String()
}

const (
	CostA = 3
	CostB = 1
)

func parsePositionRegex(matches [][]string) (Position, error) {
	if len(matches) == 0 || len(matches[0]) < 3 {
		return Position{}, fmt.Errorf("invalid match format: expected at least 3 groups, got %v", matches)
	}

	x, err1 := strconv.Atoi(matches[0][1])
	y, err2 := strconv.Atoi(matches[0][2])
	if err1 != nil {
		return Position{}, fmt.Errorf("could not convert to integer: %w", err1)
	}
	if err2 != nil {
		return Position{}, fmt.Errorf("could not convert to integer: %w", err2)
	}

	return Position{X: x, Y: y}, nil
}

func loadInput(filename string, part2 bool) ([]Game, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var games []Game
	temp := make([]string, 3)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			temp[i] = line
			i++
			if i < 3 {
				continue
			}
		}

		if i == 3 && line == "" {
			reButton := regexp.MustCompile(`Button [AB]: X\+(\d+), Y\+(\d+)`)
			rePrize := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

			match_a := reButton.FindAllStringSubmatch(temp[0], -1)
			match_b := reButton.FindAllStringSubmatch(temp[1], -1)
			match_prize := rePrize.FindAllStringSubmatch(temp[2], -1)

			button_a, err1 := parsePositionRegex(match_a)
			if err1 != nil {
				return nil, fmt.Errorf("failed to parse button A: %w", err1)
			}

			button_b, err2 := parsePositionRegex(match_b)
			if err2 != nil {
				return nil, fmt.Errorf("failed to parse button B: %w", err2)
			}

			// // DO NOT TRY THIS AT HOME
			prize_pos, err3 := parsePositionRegex(match_prize)
			if err3 != nil {
				return nil, fmt.Errorf("failed to parse prize position: %w", err3)
			}

			if part2 {
				prize_pos = prize_pos.Add(Position{10000000000000, 10000000000000})
			}
			game := Game{
				ButtonA:  button_a,
				ButtonB:  button_b,
				PrizePos: prize_pos,
			}
			games = append(games, game)
			temp = make([]string, 3)
			i = 0
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	return games, nil
}

func play(game Game, pos Position, cache map[Position]float64) float64 {
	// Found solution
	if pos == game.PrizePos {
		return 0
	}

	// No solution in this branch
	if pos.X > game.PrizePos.X || pos.Y > game.PrizePos.Y {
		return math.Inf(1)
	}

	if score, found := cache[pos]; found {
		return score
	}

	// Choose A
	a_score := play(game, pos.Add(game.ButtonA), cache) + CostA
	b_score := play(game, pos.Add(game.ButtonB), cache) + CostB

	cache[pos] = math.Min(a_score, b_score)
	return cache[pos]
}

func playWrapper(game Game) float64 {
	cache := make(map[Position]float64)
	start := Position{X: 0, Y: 0}
	return play(game, start, cache)
}

func linearAlgebraGoBrrrr(game Game) int {
	// A = span{a, b}
	// b = target vector (prize)
	/*
		Ax = b :
		[a_x, b_x]   [count_x]   [prize_x]
		[a_y, b_y] * [count_y] = [prize_y]

		inv(A) =
		1 / (a_x*b_y - b_x*a_y) *
		[b_y, -a_y]
		[-b_x, a_x]
	*/

	v_a, v_b := game.ButtonA, game.ButtonB
	b := game.PrizePos
	A := [2][2]float64{
		{float64(v_a.X), float64(v_b.X)},
		{float64(v_a.Y), float64(v_b.Y)},
	}
	// Calculate determinant
	det := A[0][0]*A[1][1] - A[0][1]*A[1][0]

	// No solution
	if det == 0 {
		return 0
	}

	// Inverse of A
	det_inv := 1 / det
	A_inv := [2][2]float64{
		{det_inv * A[1][1], -det_inv * A[1][0]},
		{-det_inv * A[0][1], det_inv * A[0][0]},
	}

	// x = inv(A)*b
	x := [2]float64{float64(b.X)*A_inv[0][0] + float64(b.Y)*A_inv[1][0], float64(b.X)*A_inv[0][1] + float64(b.Y)*A_inv[1][1]}

	// Only count positive solutions
	if x[0] < 0 || x[1] < 0 {
		return 0
	}

	// Only count integer solutions, so check if after rounding, the solutions still holds.
	sol := [2]int{int(math.Round(x[0])), int(math.Round(x[1]))}

	// Check Solution
	if !(sol[0]*v_a.X+sol[1]*v_b.X == b.X && sol[0]*v_a.Y+sol[1]*v_b.Y == b.Y) {
		return 0
	}

	// Return cost
	return sol[0]*CostA + sol[1]*CostB
}

func part1() {
	games, err := loadInput("input.txt", false)
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0.0
	for _, game := range games {
		fmt.Println(game)
		score := playWrapper(game)
		if score != math.Inf(1) {
			sum += score
		}
	}

	fmt.Println("Part I: ", sum)
}

func part2() {
	games, err := loadInput("input.txt", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, game := range games {
		sum += linearAlgebraGoBrrrr(game)
	}

	fmt.Println("Part II: ", sum)
}
func main() {
	part1()
	part2()
}
