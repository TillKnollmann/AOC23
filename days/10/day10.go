package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "10"

type Game struct {
	fields      [][]string
	startX      int
	startY      int
	limitX      int
	limitY      int
	distances   [][]int64
	tilesInLoop [][]int
}

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

type Explorer struct {
	positionX int
	positionY int
	direction Direction
}

func printDistances(game Game) {

	var output = ""

	for y := 0; y < len(game.distances); y++ {

		for x := 0; x < len(game.distances[y]); x++ {

			output += strings.Replace(fmt.Sprintf("%d", game.distances[y][x]), "0", ".", 1)
		}
		output += "\n"
	}
	fmt.Println(output)
}

func parseGame(input string) Game {

	var game Game

	for y, line := range strings.Split(input, "\n") {

		line = strings.Trim(line, "\r")

		var characters []string
		for x := 0; x < len(line); x++ {

			characters = append(characters, string(line[x]))
			if string(line[x]) == "S" {

				game.startY = y
				game.startX = x
			}
		}

		game.fields = append(game.fields, characters)
		game.distances = append(game.distances, make([]int64, len(characters)))
		game.tilesInLoop = append(game.tilesInLoop, make([]int, len(characters)))
	}

	game.limitY = len(game.fields)
	game.limitX = len(game.fields[0])

	game.tilesInLoop[game.startY][game.startX] = 2

	return game
}

func calculateDistances(game *Game) {

	foundLoop := calculateDistancesWithExplorer(game, Explorer{positionX: game.startX, positionY: game.startY, direction: North}, 0, false)
	foundLoop = calculateDistancesWithExplorer(game, Explorer{positionX: game.startX, positionY: game.startY, direction: East}, 0, foundLoop)
	foundLoop = calculateDistancesWithExplorer(game, Explorer{positionX: game.startX, positionY: game.startY, direction: South}, 0, foundLoop)
	foundLoop = calculateDistancesWithExplorer(game, Explorer{positionX: game.startX, positionY: game.startY, direction: West}, 0, foundLoop)

}

func calculateDistancesWithExplorer(game *Game, explorer Explorer, distance int64, loopAlreadyFound bool) bool {

	moveExplorer(&explorer)

	newDistance := distance + 1

	if explorer.positionX < 0 || explorer.positionX >= game.limitX || explorer.positionY < 0 || explorer.positionY >= game.limitY {

		return false
	}

	switch game.fields[explorer.positionY][explorer.positionX] {
	case ".":
		return false
	case "S":
		return true
	case "|":
		if !(explorer.direction == South || explorer.direction == North) {
			return false
		}
	case "-":
		if !(explorer.direction == West || explorer.direction == East) {
			return false
		}
	case "L":
		if !(explorer.direction == South || explorer.direction == West) {
			return false
		}
		if explorer.direction == South {
			explorer.direction = East
		} else if explorer.direction == West {
			explorer.direction = North
		}
	case "J":
		if !(explorer.direction == South || explorer.direction == East) {
			return false
		}
		if explorer.direction == South {
			explorer.direction = West
		} else if explorer.direction == East {
			explorer.direction = North
		}
	case "7":
		if !(explorer.direction == North || explorer.direction == East) {
			return false
		}
		if explorer.direction == North {
			explorer.direction = West
		} else if explorer.direction == East {
			explorer.direction = South
		}
	case "F":
		if !(explorer.direction == North || explorer.direction == West) {
			return false
		}
		if explorer.direction == North {
			explorer.direction = East
		} else if explorer.direction == West {
			explorer.direction = South
		}
	}

	if game.distances[explorer.positionY][explorer.positionX] == 0 || game.distances[explorer.positionY][explorer.positionX] > newDistance {
		game.distances[explorer.positionY][explorer.positionX] = newDistance
	}
	loopFound := calculateDistancesWithExplorer(game, explorer, newDistance, loopAlreadyFound)

	if !(loopFound) {

		game.distances[explorer.positionY][explorer.positionX] = 0
	} else {

		// 2 means there is a pipe of the loop
		game.tilesInLoop[explorer.positionY][explorer.positionX] = 2
	}

	if loopAlreadyFound {

		// mark tiles
		switch game.fields[explorer.positionY][explorer.positionX] {
		case "|":
			if explorer.direction == North {
				floodTiles(game, explorer.positionX+1, explorer.positionY)
			} else if explorer.direction == South {
				floodTiles(game, explorer.positionX-1, explorer.positionY)
			}
		case "-":
			if explorer.direction == East {
				floodTiles(game, explorer.positionX, explorer.positionY+1)
			} else if explorer.direction == West {
				floodTiles(game, explorer.positionX, explorer.positionY-1)
			}
		case "L":
			if explorer.direction == East {
				floodTiles(game, explorer.positionX-1, explorer.positionY)
				floodTiles(game, explorer.positionX, explorer.positionY+1)
				floodTiles(game, explorer.positionX-1, explorer.positionY+1)
			}
		case "J":
			if explorer.direction == North {
				floodTiles(game, explorer.positionX+1, explorer.positionY)
				floodTiles(game, explorer.positionX, explorer.positionY+1)
				floodTiles(game, explorer.positionX+1, explorer.positionY+1)
			}
		case "7":
			if explorer.direction == West {
				floodTiles(game, explorer.positionX, explorer.positionY-1)
				floodTiles(game, explorer.positionX+1, explorer.positionY)
				floodTiles(game, explorer.positionX+1, explorer.positionY-1)
			}
		case "F":
			if explorer.direction == South {
				floodTiles(game, explorer.positionX-1, explorer.positionY)
				floodTiles(game, explorer.positionX, explorer.positionY-1)
				floodTiles(game, explorer.positionX-1, explorer.positionY-1)
			}
		}

	}

	return loopFound
}

func floodTiles(game *Game, positionX int, positionY int) {

	if !isInBounds(*game, positionX, positionY) {
		return
	}

	if game.tilesInLoop[positionY][positionX] != 0 {
		return
	}

	// marks a tile
	game.tilesInLoop[positionY][positionX] = 1

	floodTiles(game, positionX-1, positionY)
	floodTiles(game, positionX+1, positionY)
	floodTiles(game, positionX, positionY-1)
	floodTiles(game, positionX, positionY+1)
}

func countMarkedTiles(game Game) (int64, int64) {

	sumMarked := int64(0)
	sumUnmarked := int64(0)

	for y := 0; y < len(game.tilesInLoop); y++ {

		for x := 0; x < len(game.tilesInLoop[y]); x++ {

			if game.tilesInLoop[y][x] == 1 {

				sumMarked++
			} else if game.tilesInLoop[y][x] == 0 {

				sumUnmarked++
			}
		}
	}

	return sumMarked, sumUnmarked
}

func isInBounds(game Game, positionX int, positionY int) bool {

	return 0 <= positionX && positionX < game.limitX && 0 <= positionY && positionY < game.limitY
}

func moveExplorer(explorer *Explorer) {

	if explorer.direction == North {

		explorer.positionY--
	} else if explorer.direction == East {

		explorer.positionX++
	} else if explorer.direction == South {

		explorer.positionY++
	} else if explorer.direction == West {

		explorer.positionX--
	}
}

func getMax(array [][]int64) int64 {

	maxValue := int64(0)

	for y := 0; y < len(array); y++ {
		for x := 0; x < len(array[y]); x++ {
			maxValue = max(maxValue, array[y][x])
		}
	}

	return maxValue
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	calculateDistances(&game)

	return strconv.FormatInt(getMax(game.distances), 10)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	calculateDistances(&game)

	optionA, optionB := countMarkedTiles(game)

	// just guess that it is the minimum lol
	return strconv.FormatInt(min(optionB, optionA), 10)
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
