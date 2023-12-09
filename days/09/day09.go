package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DAY = "09"

type Game struct {
	levels [][]int64
}

func (game Game) String() string {

	return fmt.Sprintf("Game(levels=%#v)", game.levels)
}

func parseGame(input string) Game {

	numberRe := regexp.MustCompile(`-?\d+`)

	var levelZero []int64

	levelZero = append(levelZero, -1)

	for _, number := range numberRe.FindAllString(input, -1) {

		levelZero = append(levelZero, stringToNumber(number))
	}

	levelZero = append(levelZero, -1)

	var game Game
	game.levels = append(game.levels, levelZero)

	return game
}

func initializeNewLevel(game *Game) {

	newLevel := make([]int64, len(game.levels[len(game.levels)-1])-1)

	for index := 0; index < len(newLevel); index++ {

		newLevel[index] = -1
	}

	game.levels = append(game.levels, newLevel)
}

func calculateFromPrevious(game *Game, level int, position int) {

	game.levels[level][position] = game.levels[level-1][position+1] - game.levels[level-1][position]
}

func calculateLevelFromPrevious(game *Game) {

	initializeNewLevel(game)

	level := len(game.levels) - 1

	for position := 1; position < len(game.levels[level])-1; position++ {

		calculateFromPrevious(game, level, position)
	}
}

func calculateTopDown(game *Game) {

	for {

		isDone := true
		for index := 1; index < len(game.levels[len(game.levels)-1])-1; index++ {

			if game.levels[len(game.levels)-1][index] != 0 {

				isDone = false
				break
			}
		}

		if isDone {

			game.levels[len(game.levels)-1][len(game.levels[len(game.levels)-1])-1] = 0
			game.levels[len(game.levels)-1][0] = 0
			break
		}

		calculateLevelFromPrevious(game)
	}
}

func calculateBottomUp(game *Game) {

	currentLevel := len(game.levels) - 2
	currentPositionStart := 0
	currentPositionEnd := len(game.levels[currentLevel]) - 1

	for {
		if currentLevel == -1 {
			break
		}

		game.levels[currentLevel][currentPositionEnd] = game.levels[currentLevel][currentPositionEnd-1] + game.levels[currentLevel+1][currentPositionEnd-1]
		game.levels[currentLevel][currentPositionStart] = game.levels[currentLevel][currentPositionStart+1] - game.levels[currentLevel+1][currentPositionStart]

		currentLevel--
		currentPositionEnd++
	}
}

func calculateGame(game *Game) {

	calculateTopDown(game)
	calculateBottomUp(game)
}

func Part1(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	var sum int64 = int64(0)

	for _, line := range lines {

		game := parseGame(line)
		calculateGame(&game)
		sum += game.levels[0][len(game.levels[0])-1]
	}

	return strconv.FormatInt(sum, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	var sum int64 = int64(0)

	for _, line := range lines {

		game := parseGame(line)
		calculateGame(&game)
		sum += game.levels[0][0]
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
