package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "18"

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

type Game struct {
	commands []DigCommand
	fields   [][]string
	start    Point
}

type DigCommand struct {
	direction Direction
	length    int
	color     string
}

type Point struct {
	positionX int
	positionY int
}

func parseGame(input string) Game {

	var game Game

	game.commands = parseCommands(input)

	minX := 0
	maxX := 0
	minY := 0
	maxY := 0

	current := Point{
		positionX: 0,
		positionY: 0,
	}

	for _, command := range game.commands {

		switch command.direction {
		case North:
			current.positionY -= command.length
		case East:
			current.positionX += command.length
		case South:
			current.positionY += command.length
		case West:
			current.positionX -= command.length
		}
		minX = min(minX, current.positionX)
		maxX = max(maxX, current.positionX)
		minY = min(minY, current.positionY)
		maxY = max(maxY, current.positionY)
	}

	game.start = Point{
		positionX: -1 * minX,
		positionY: -1 * minY,
	}

	sizeX := -1*minX + maxX + 1
	sizeY := -1*minY + maxY + 1

	for y := 0; y < sizeY; y++ {

		var row []string
		for x := 0; x < sizeX; x++ {

			row = append(row, "_")
		}
		game.fields = append(game.fields, row)
	}

	return game
}

func parseCommands(input string) []DigCommand {

	var commands []DigCommand

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	for _, line := range lines {

		parts := strings.Split(line, " ")

		direction := North
		switch parts[0] {
		case `R`:
			direction = East
		case `D`:
			direction = South
		case `L`:
			direction = West
		}

		commands = append(commands, DigCommand{
			direction: direction,
			length:    int(stringToNumber(parts[1])),
			color:     parts[2],
		})
	}

	return commands
}

func digBorder(game *Game) {

	current := game.start
	game.fields[current.positionY][current.positionX] = "#"

	for _, command := range game.commands {

		for i := 0; i < command.length; i++ {

			switch command.direction {
			case North:
				current.positionY--
			case East:
				current.positionX++
			case South:
				current.positionY++
			case West:
				current.positionX--
			}

			game.fields[current.positionY][current.positionX] = "#"
		}
	}
}

func digInterior(game *Game) {

	current := game.start

	for _, command := range game.commands {

		for i := 0; i < command.length; i++ {

			switch command.direction {
			case North:
				current.positionY--
			case East:
				current.positionX++
			case South:
				current.positionY++
			case West:
				current.positionX--
			}

			var pointToTheRight = Point{
				positionX: current.positionX,
				positionY: current.positionY,
			}

			switch command.direction {
			case North:
				pointToTheRight.positionX++
			case East:
				pointToTheRight.positionY++
			case South:
				pointToTheRight.positionX--
			case West:
				pointToTheRight.positionY--
			}

			floodInterior(game, pointToTheRight)
		}
	}
}

func floodInterior(game *Game, point Point) {

	if point.positionX < 0 || point.positionX >= len(game.fields[0]) || point.positionY < 0 || point.positionY >= len(game.fields) {

		return
	}

	if game.fields[point.positionY][point.positionX] == "#" {

		return
	}

	game.fields[point.positionY][point.positionX] = "#"

	floodInterior(game, Point{
		positionX: point.positionX + 1,
		positionY: point.positionY,
	})

	floodInterior(game, Point{
		positionX: point.positionX - 1,
		positionY: point.positionY,
	})

	floodInterior(game, Point{
		positionX: point.positionX,
		positionY: point.positionY + 1,
	})

	floodInterior(game, Point{
		positionX: point.positionX,
		positionY: point.positionY - 1,
	})
}

func getFilledCount(game Game) int {

	count := 0

	for y := 0; y < len(game.fields); y++ {

		for x := 0; x < len(game.fields[y]); x++ {

			if game.fields[y][x] == "#" {

				count++
			}
		}
	}

	return count
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	digBorder(&game)

	digInterior(&game)

	return strconv.Itoa(getFilledCount(game))
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
