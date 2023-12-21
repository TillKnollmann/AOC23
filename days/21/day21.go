package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DAY = "21"

type Game struct {
	nodesById       map[int]*Node
	nodesByPosition [][]*Node
	start           *Node
}

type Node struct {
	x           int
	y           int
	character   string
	id          int
	adjacentIds []int
}

type AdjacencyMatrix struct {
	items [][]int
}

func calculateAdjacencyMatrix(game *Game) AdjacencyMatrix {

	for y := 0; y < len(game.nodesByPosition); y++ {
		for x := 0; x < len(game.nodesByPosition[y]); x++ {

			// check each node that is "S" or "#"
			if game.nodesByPosition[y][x].character != "." {

				// check all positions around the node

				for deltaX := -1; deltaX <= 1; deltaX += 2 {

					if x+deltaX >= 0 && x+deltaX < len(game.nodesByPosition[y]) && game.nodesByPosition[y][x+deltaX].character != "." {

						game.nodesByPosition[y][x].adjacentIds = append(game.nodesByPosition[y][x].adjacentIds, game.nodesByPosition[y][x+deltaX].id)
					}
				}

				for deltaY := -1; deltaY <= 1; deltaY += 2 {

					// the position is adjacent iff it is within range and either "S" or "#"
					if y+deltaY >= 0 && y+deltaY < len(game.nodesByPosition) && game.nodesByPosition[y+deltaY][x].character != "." {

						game.nodesByPosition[y][x].adjacentIds = append(game.nodesByPosition[y][x].adjacentIds, game.nodesByPosition[y+deltaY][x].id)
					}
				}
			}
		}
	}

	var matrix AdjacencyMatrix

	for i := 0; i < len(game.nodesById); i++ {

		nodeOne, _ := game.nodesById[i]

		var row []int
		for j := 0; j < len(game.nodesById); j++ {

			nodeTwo, _ := game.nodesById[j]

			if slices.Contains(nodeOne.adjacentIds, nodeTwo.id) {

				row = append(row, 1)

			} else {

				row = append(row, 0)

			}
		}

		matrix.items = append(matrix.items, row)
	}

	return matrix
}

func parseGame(input string, maxSteps int) Game {

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	// first identify start point
	startX := 0
	startY := 0
	found := false

	for y := 0; y < len(lines) && !found; y++ {
		for x := 0; x < len(lines[0]) && !found; x++ {
			if string(lines[y][x]) == "S" {
				startY = y
				startX = x
				found = true
			}
		}
	}

	// get frame of game to consider
	offsetX := max(0, startX-maxSteps-1)
	offsetY := max(0, startY-maxSteps-1)

	var game Game

	game.nodesById = make(map[int]*Node)

	id := 0

	for y := offsetY; y < min(len(lines), startY+offsetY); y++ {

		var row []*Node
		for x := offsetX; x < min(len(lines[y]), startX+offsetX); x++ {

			node := Node{
				x:           x - offsetX,
				y:           y - offsetY,
				character:   string(lines[y][x]),
				id:          id,
				adjacentIds: []int{},
			}

			game.nodesById[id] = &node
			row = append(row, &node)

			if node.character == "S" {

				game.start = &node
			}

			id++
		}

		game.nodesByPosition = append(game.nodesByPosition, row)
	}

	return game
}

func power(matrix [][]int, x int) [][]int {
	n := len(matrix)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
		result[i][i] = 1
	}
	for ; x > 0; x >>= 1 {
		if x&1 == 1 {
			result = multiply(result, matrix)
		}
		matrix = multiply(matrix, matrix)
	}
	return result
}

func multiply(a, b [][]int) [][]int {
	n := len(a)
	c := make([][]int, n)
	for i := range c {
		c[i] = make([]int, n)
		for j := range c[i] {
			for k := range a[i] {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func parallelPower(matrix [][]int, x int) [][]int {
	n := len(matrix)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
		result[i][i] = 1
	}
	done := make(chan bool)
	for ; x > 0; x >>= 1 {
		if x&1 == 1 {
			result = parallelMultiply(result, matrix, done)
		}
		matrix = parallelMultiply(matrix, matrix, done)
	}
	return result
}

func parallelMultiply(a, b [][]int, done chan bool) [][]int {
	n := len(a)
	c := make([][]int, n)
	for i := range c {
		c[i] = make([]int, n)
		for j := range c[i] {
			c[i][j] = 0
			go func(i, j int) {
				for k := range a[i] {
					c[i][j] += a[i][k] * b[k][j]
				}
				done <- true
			}(i, j)
		}
	}
	for i := 0; i < n*n; i++ {
		<-done
	}
	return c
}

func Part1(input string) string {

	content := GetContent(input)

	maxSteps := 64
	powerParam := 6

	game := parseGame(content, maxSteps)

	matrix := calculateAdjacencyMatrix(&game)

	// since 2^6 = 64, six time square the adjacency matrix to get all positions in distance 64 in row of S
	matrix.items = parallelPower(matrix.items, powerParam)

	// the answer is given by the non-zero entries in row of S
	count := 0
	for x := 0; x < len(matrix.items[game.start.id]); x++ {

		if matrix.items[game.start.id][x] != 0 {
			count++
		}
	}

	return strconv.Itoa(count)
}

func Part2(input string) string {

	content := GetContent(input)

	return string(content[0])
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
