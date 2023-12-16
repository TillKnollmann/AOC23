package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DAY = "16"

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

type Tile struct {
	character     string
	outgoingBeams []Direction
}

type Beam struct {
	direction Direction
	positionX int
	positionY int
}

type Game struct {
	tiles      [][]Tile
	limitX     int
	limitY     int
	exitPoints []Point
}

type Point struct {
	positionX int
	positionY int
}

func parseGame(input string) Game {

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	var game Game

	for y := 0; y < len(lines); y++ {

		var row []Tile
		for x := 0; x < len(lines[y]); x++ {

			row = append(row, Tile{character: string(lines[y][x])})
		}

		game.tiles = append(game.tiles, row)
	}

	game.limitY = len(lines)
	game.limitX = len(lines[0])

	return game
}

func calculateBeam(beam Beam, game *Game) {

	// catch beams off playground
	if beam.positionY < 0 || beam.positionY >= game.limitY || beam.positionX < 0 || beam.positionX >= game.limitX {

		exitPoint := Point{
			positionX: beam.positionX,
			positionY: beam.positionY,
		}

		if !slices.Contains(game.exitPoints, exitPoint) {

			game.exitPoints = append(game.exitPoints, exitPoint)
		}
		return
	}

	// calculate resulting beams
	var newBeams []Beam

	switch game.tiles[beam.positionY][beam.positionX].character {
	case ".":
		newBeams = append(newBeams, Beam{
			direction: beam.direction,
			positionX: beam.positionX,
			positionY: beam.positionY,
		})
		break
	case "\\":
		var newDirection Direction
		switch beam.direction {
		case North:
			newDirection = West
		case East:
			newDirection = South
		case South:
			newDirection = East
		case West:
			newDirection = North
		}
		newBeams = append(newBeams, Beam{
			direction: newDirection,
			positionX: beam.positionX,
			positionY: beam.positionY,
		})
	case "/":
		var newDirection Direction
		switch beam.direction {
		case North:
			newDirection = East
		case East:
			newDirection = North
		case South:
			newDirection = West
		case West:
			newDirection = South
		}
		newBeams = append(newBeams, Beam{
			direction: newDirection,
			positionX: beam.positionX,
			positionY: beam.positionY,
		})
	case "-":
		if beam.direction == West || beam.direction == East {
			newBeams = append(newBeams, Beam{
				direction: beam.direction,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
		} else {
			newBeams = append(newBeams, Beam{
				direction: West,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
			newBeams = append(newBeams, Beam{
				direction: East,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
		}
	case "|":
		if beam.direction == North || beam.direction == South {
			newBeams = append(newBeams, Beam{
				direction: beam.direction,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
		} else {
			newBeams = append(newBeams, Beam{
				direction: North,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
			newBeams = append(newBeams, Beam{
				direction: South,
				positionX: beam.positionX,
				positionY: beam.positionY,
			})
		}
	}

	// remove resulting beams that already existed on the tile
	var prunedNewBeams []Beam
	for _, newBeam := range newBeams {

		if !slices.Contains(game.tiles[beam.positionY][beam.positionX].outgoingBeams, newBeam.direction) {

			prunedNewBeams = append(prunedNewBeams, newBeam)
			game.tiles[beam.positionY][beam.positionX].outgoingBeams = append(game.tiles[beam.positionY][beam.positionX].outgoingBeams, newBeam.direction)
		}
	}

	// calculated moved beams
	var movedNewBeams []Beam
	for _, prunedBeam := range prunedNewBeams {

		moveBeam(&prunedBeam)
		movedNewBeams = append(movedNewBeams, prunedBeam)
	}

	// recurse
	for _, movedBeam := range movedNewBeams {

		calculateBeam(movedBeam, game)
	}
}

func moveBeam(beam *Beam) {

	newX, newY := getPosition(beam.positionX, beam.positionY, beam.direction)
	beam.positionX = newX
	beam.positionY = newY
}

func getPosition(positionX int, positionY int, direction Direction) (int, int) {

	switch direction {
	case North:
		return positionX, positionY - 1
	case East:
		return positionX + 1, positionY
	case South:
		return positionX, positionY + 1
	case West:
		return positionX - 1, positionY
	}

	return 0, 0
}

func getEnergizedTiles(game Game) int {

	sum := 0

	for y := 0; y < game.limitY; y++ {

		for x := 0; x < game.limitX; x++ {

			if isTileEnergized(game.tiles[y][x]) {

				sum++
			}
		}
	}

	return sum
}

func isTileEnergized(tile Tile) bool {

	return len(tile.outgoingBeams) > 0
}

func flushTiles(game *Game) {

	for y := 0; y < game.limitY; y++ {

		for x := 0; x < game.limitX; x++ {

			game.tiles[y][x].outgoingBeams = slices.DeleteFunc(game.tiles[y][x].outgoingBeams, func(direction Direction) bool {
				return true
			})
		}
	}
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)
	calculateBeam(Beam{
		direction: East,
		positionX: 0,
		positionY: 0,
	}, &game)

	result := getEnergizedTiles(game)

	return strconv.Itoa(result)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	// get all start beams (initially outside the field)
	var startBeams []Beam
	for y := 0; y < game.limitY; y++ {

		startBeams = append(startBeams, Beam{
			direction: East,
			positionX: -1,
			positionY: y,
		})
		startBeams = append(startBeams, Beam{
			direction: West,
			positionX: game.limitX,
			positionY: y,
		})
	}

	for x := 0; x < game.limitX; x++ {

		startBeams = append(startBeams, Beam{
			direction: South,
			positionX: x,
			positionY: -1,
		})
		startBeams = append(startBeams, Beam{
			direction: North,
			positionX: x,
			positionY: game.limitY,
		})
	}

	maxEnergized := 0

	for len(startBeams) > 0 {

		startBeam := startBeams[0]
		startBeams = slices.Delete(startBeams, 0, 1)

		// startBeam enters the field
		moveBeam(&startBeam)

		// calculate energized tiles
		calculateBeam(startBeam, &game)
		maxEnergized = max(maxEnergized, getEnergizedTiles(game))

		// reset field
		flushTiles(&game)

		// remove all start beams that start at a position already captured by an exit position
		startBeams = slices.DeleteFunc(startBeams, func(beam Beam) bool {

			for _, point := range game.exitPoints {

				if point.positionX == beam.positionX && point.positionY == beam.positionY {

					return true
				}
			}
			return false
		})
	}

	return strconv.Itoa(maxEnergized)

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
