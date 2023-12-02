package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Color int

const (
	Red Color = iota + 1
	Blue
	Green
	NotAColor
)

type Game struct {
	id    int
	draws []Draw
}

type Draw struct {
	ballCounts []BallCount
}

type BallCount struct {
	color Color
	count int
}

func (b Game) String() string {

	return fmt.Sprintf("Game(id=%d, draws=%s)", b.id, b.draws)
}

func (b Draw) String() string {

	return fmt.Sprintf("Draw(ballCounts=%s)", b.ballCounts)
}

func (b BallCount) String() string {

	return fmt.Sprintf("BallCount(count=%d, color=%d)", b.count, b.color)
}

func parseGame(s string) Game {

	s = strings.TrimSpace(s)

	split := strings.Split(s, ":")

	gameString := split[0]
	drawsString := split[1]

	re := regexp.MustCompile(`\d+`)

	var game Game

	gameId, _ := strconv.ParseInt(re.FindString(gameString), 10, 8)
	game.id = int(gameId)

	for _, drawString := range strings.Split(drawsString, ";") {

		game.draws = append(game.draws, parseDraw(drawString))
	}

	return game
}

func parseDraw(s string) Draw {

	s = strings.TrimSpace(s)

	ballCountStrings := strings.Split(s, ",")

	var draw Draw

	for _, ballCountString := range ballCountStrings {

		draw.ballCounts = append(draw.ballCounts, parseBallCount(ballCountString))
	}

	return draw
}

func parseBallCount(s string) BallCount {

	s = strings.TrimSpace(s)

	countString := strings.Split(s, " ")[0]
	colorString := strings.Split(s, " ")[1]

	count, _ := strconv.ParseInt(countString, 10, 8)
	countInt := int(count)
	color := parseColor(colorString)

	var ballCount BallCount

	ballCount.count = countInt
	ballCount.color = color

	return ballCount
}

func parseColor(s string) Color {

	s = strings.TrimSpace(s)

	var result Color

	switch s {

	case "red":
		result = Red
	case "green":
		result = Green
	case "blue":
		result = Blue
	default:
		result = NotAColor

	}

	return result
}

func isGamePossible(game Game, limitRed int, limitGreen int, limitBlue int) bool {

	for _, draw := range game.draws {

		if !isDrawPossible(draw, limitRed, limitGreen, limitBlue) {

			return false
		}
	}

	return true
}

func isDrawPossible(draw Draw, limitRed int, limitGreen int, limitBlue int) bool {

	for _, ballCount := range draw.ballCounts {

		var limit int

		switch ballCount.color {
		case Red:
			limit = limitRed
		case Green:
			limit = limitGreen
		case Blue:
			limit = limitBlue
		}

		if ballCount.count > limit {

			return false
		}
	}

	return true
}

func getGamePower(game Game) int64 {

	var maxRed int64 = 0
	var maxGreen int64 = 0
	var maxBlue int64 = 0

	for _, draw := range game.draws {

		for _, ballCount := range draw.ballCounts {

			switch ballCount.color {
			case Red:
				maxRed = max(maxRed, int64(ballCount.count))
			case Green:
				maxGreen = max(maxGreen, int64(ballCount.count))
			case Blue:
				maxBlue = max(maxBlue, int64(ballCount.count))
			}
		}
	}

	if maxRed == 0 || maxGreen == 0 || maxBlue == 0 {

		return 0
	}

	return maxRed * maxGreen * maxBlue
}

func Part1(input string) string {

	limitRed := 12
	limitGreen := 13
	limitBlue := 14

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	sum := 0

	for _, line := range lines {

		game := parseGame(line)
		if isGamePossible(game, limitRed, limitGreen, limitBlue) {

			sum += game.id
		}
	}

	return strconv.Itoa(sum)
}

func Part2(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	var sum int64 = 0

	for _, line := range lines {

		game := parseGame(line)
		sum += getGamePower(game)
	}

	return strconv.FormatInt(sum, 10)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(Part1("input/02/in01.txt"))
	fmt.Println(Part2("input/02/in01.txt"))
}
