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

			// check each node that is "S" or "."
			if game.nodesByPosition[y][x].character != "#" {

				// check all positions around the node

				for deltaX := -1; deltaX <= 1; deltaX += 2 {

					if x+deltaX >= 0 && x+deltaX < len(game.nodesByPosition[y]) && game.nodesByPosition[y][x+deltaX].character != "#" {

						game.nodesByPosition[y][x].adjacentIds = append(game.nodesByPosition[y][x].adjacentIds, game.nodesByPosition[y][x+deltaX].id)
					}
				}

				for deltaY := -1; deltaY <= 1; deltaY += 2 {

					// the position is adjacent iff it is within range and either "S" or "."
					if y+deltaY >= 0 && y+deltaY < len(game.nodesByPosition) && game.nodesByPosition[y+deltaY][x].character != "#" {

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

func parseGame(input string) Game {

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	var game Game

	game.nodesById = make(map[int]*Node)

	id := 0

	for y := 0; y < len(lines); y++ {

		var row []*Node
		for x := 0; x < len(lines[y]); x++ {

			node := Node{
				x:           x,
				y:           y,
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

func findNodes(matrix [][]int, id int, depth int) []int {

	nodes := []int{id}
	for i := 0; i < depth; i++ {
		newNodes := []int{}
		for _, node := range nodes {
			for j, edge := range matrix[node] {
				if edge == 1 {
					newNodes = append(newNodes, j)
				}
			}
		}
		nodes = removeDuplicates(newNodes)
	}
	return nodes
}

func removeDuplicates(slice []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for _, v := range slice {
		if encountered[v] == true {
			continue
		} else {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func distance(game Game, nodeAId int, nodeBId int) int {

	nodeA := game.nodesById[nodeAId]
	nodeB := game.nodesById[nodeBId]

	return max(nodeA.x, nodeB.x) - min(nodeA.x, nodeB.x) + (max(nodeA.y, nodeB.y) - min(nodeA.y, nodeB.y))
}

func Part1(input string, maxSteps int) string {

	content := GetContent(input)

	game := parseGame(content)

	matrix := calculateAdjacencyMatrix(&game)

	// calculate reachable nodes
	count := len(findNodes(matrix.items, game.start.id, maxSteps))

	return strconv.Itoa(count)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	matrix := calculateAdjacencyMatrix(&game)

	// total steps: 26501365
	// diameter of a square: 65 + 1 + 65 = 131
	// number of squares right to the initial one in one direction: (26501365 - 65) / 131 = 202300
	// number of filled squares around the initial one ending up odd: (202300 + 1)^2 = 40925694601
	// number of filled squares around the initial one ending up even: (202300)^2 = 40925290000
	// number corner squares to add (always 4 corner squares add up): 202300
	// number of corner squares to remove: 202300 + 1 = 202301

	const squaresOdd = 40925694601
	const squaresEven = 40925290000
	const squaresToAdd = 202300
	const squaresToRemove = 202301

	// get nodes that are reachable after an even number of steps and an odd number of steps
	nodesEven := findNodes(matrix.items, game.start.id, 132)
	nodesOdd := findNodes(matrix.items, game.start.id, 131)

	// calculate reachable tiles in odd/even number of steps
	tilesReachableEven := len(nodesEven)
	tilesReachableOdd := len(nodesOdd)

	// corners are nodes in a manhattan distance of more than diameter of square from center
	tilesCornerEven := 0

	for _, nodeId := range nodesEven {

		if distance(game, nodeId, game.start.id) > 65 {

			tilesCornerEven++
		}
	}

	tilesCornerOdd := 0

	for _, nodeId := range nodesOdd {

		if distance(game, nodeId, game.start.id) > 65 {

			tilesCornerOdd++
		}
	}

	// all corners to add are of even squares, all corners to remove of odd squares
	return strconv.Itoa(squaresOdd*tilesReachableOdd + squaresEven*tilesReachableEven + squaresToAdd*tilesCornerEven - squaresToRemove*tilesCornerOdd)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY), 64)))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
