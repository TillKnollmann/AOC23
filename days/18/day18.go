package main

import (
	"fmt"
	"log"
	"math"
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
	vertices []Point
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

func parseGame(input string, isPartTwo bool) Game {

	var game Game

	game.commands = parseCommands(input)

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

func calculateVertices(game *Game) {

	current := Point{
		positionX: 0,
		positionY: 0,
	}

	game.vertices = append(game.vertices, current)

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

		game.vertices = append(game.vertices, current)
	}
}

func parseColorsToCommands(game *Game) {

	var newCommands []DigCommand

	for _, command := range game.commands {

		direction, length := parseColor(command.color)

		newCommands = append(newCommands, DigCommand{
			direction: direction,
			length:    length,
			color:     "",
		})
	}

	game.commands = newCommands
}

func parseColor(color string) (Direction, int) {

	length, _ := strconv.ParseInt(string(color[2:len(color)-2]), 16, 64)

	switch string(color[len(color)-2 : len(color)-1]) {
	case "0":
		return East, int(length)
	case "1":
		return South, int(length)
	case "2":
		return West, int(length)
	case "3":
		return North, int(length)
	}

	return North, 0
}

func getFilledCount(game Game) int {

	return int(ShoelaceFormula(game.vertices)) + int(Perimeter(game.vertices)/2) - 1
}

func ShoelaceFormula(points []Point) float64 {
	var sum float64
	n := len(points)
	for i := 0; i < n-1; i++ {
		sum += float64(points[i].positionX*points[i+1].positionY - points[i+1].positionX*points[i].positionY)
	}
	return math.Abs(sum / 2)
}

func Perimeter(points []Point) float64 {
	var sum float64
	n := len(points)
	for i := 0; i < n-1; i++ {
		sum += math.Abs(float64(points[i+1].positionX-points[i].positionX)) + math.Abs(float64(points[i+1].positionY-points[i].positionY))
	}
	return sum + 4
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content, false)

	calculateVertices(&game)

	return strconv.Itoa(getFilledCount(game))
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content, true)

	parseColorsToCommands(&game)

	calculateVertices(&game)

	return strconv.Itoa(getFilledCount(game))
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
