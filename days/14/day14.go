package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DAY = "14"

type Game struct {
	columns []string
}

func parseGame(input string) Game {

	var game Game

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	game.columns = lines
	rotateGameClockwise(&game)

	slices.Reverse(game.columns)

	return game
}

func rotateGameClockwise(game *Game) {

	var newColumns []string

	for y := len(game.columns[0]) - 1; y >= 0; y-- {

		newColumn := ""

		for x := 0; x < len(game.columns); x++ {

			newColumn += string(game.columns[x][y])
		}

		newColumns = append(newColumns, newColumn)
	}

	game.columns = newColumns
}

func getRockChunksPerColumn(game Game) [][]int {

	var rockChunksPerColumn [][]int

	for _, column := range game.columns {

		var rockChunks []int
		currentRockCount := 0

		for i := 0; i < len(column); i++ {

			if string(column[i]) == "O" {

				currentRockCount++

			} else if string(column[i]) == "#" {

				rockChunks = append(rockChunks, currentRockCount)
				currentRockCount = 0
			}
		}

		rockChunks = append(rockChunks, currentRockCount)

		rockChunksPerColumn = append(rockChunksPerColumn, rockChunks)
	}

	return rockChunksPerColumn
}

func tiltGameNorth(game *Game) {

	rocksInColumn := getRockChunksPerColumn(*game)

	var tiltedColumns []string

	for columnIndex, column := range game.columns {

		tiltedColumn := ""

		rockChunk := rocksInColumn[columnIndex]
		rockChunkIndex := 0
		placedRocks := 0

		for i := 0; i < len(column); i++ {

			if rockChunkIndex == len(rockChunk) {

				tiltedColumn += column[i:]
				break
			}

			if string(column[i]) == "." || string(column[i]) == "O" {

				if placedRocks < rockChunk[rockChunkIndex] {

					tiltedColumn += "O"
					placedRocks++
				} else {

					tiltedColumn += "."
				}
			} else if string(column[i]) == "#" {

				tiltedColumn += "#"
				rockChunkIndex++
				placedRocks = 0
			}
		}

		tiltedColumns = append(tiltedColumns, tiltedColumn)
	}

	game.columns = tiltedColumns
}

func tiltCycle(game *Game) {

	for i := 0; i < 4; i++ {

		tiltGameNorth(game)
		rotateGameClockwise(game)
	}
}

func getTotalLoad(game Game) int {

	sum := 0

	for _, column := range game.columns {

		for y := 0; y < len(column); y++ {

			if string(column[y]) == "O" {

				score := len(column) - y
				sum += score
			}
		}
	}

	return sum
}

func checkForLoop(loads []int, games []string, cycleIndices []int) (bool, int) {

	for index := len(loads) - 2; index >= 0; index-- {

		if loads[index] == loads[len(loads)-1] && games[index] == games[len(loads)-1] {

			return true, cycleIndices[len(cycleIndices)-1] - cycleIndices[index]
		}
	}

	return false, 0
}

func stringify(game Game) string {

	return strings.Join(game.columns, "\n")
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	tiltGameNorth(&game)

	sum := getTotalLoad(game)

	return strconv.Itoa(sum)
}

func Part2(input string) string {

	content := GetContent(input)

	maxCycles := 1000000000

	game := parseGame(content)

	var loads []int
	var cycleIndices []int
	var games []string

	for i := 0; i < maxCycles; i++ {

		tiltCycle(&game)

		loads = append(loads, getTotalLoad(game))
		cycleIndices = append(cycleIndices, i)
		games = append(games, stringify(game))

		foundLoop, deltaIndex := checkForLoop(loads, games, cycleIndices)

		if foundLoop {

			increase := int(math.Floor(float64((maxCycles-1-i)/deltaIndex))) * deltaIndex
			i += increase
		}
	}

	sum := getTotalLoad(game)

	return strconv.Itoa(sum)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func stringToNumber(s string) int64 {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return number
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
