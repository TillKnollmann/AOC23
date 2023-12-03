package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DAY = "03"

func getPartNumbers(board []string) []int {

	limitY := len(board) - 1
	limitX := len(board[0]) - 1

	var partNumbers []int

	re := regexp.MustCompile(`\d+`)

	for indexY, line := range board {

		numberLocations := re.FindAllStringIndex(line, -1)

		for _, numberLocation := range numberLocations {

			var adjacentFields []string

			bounceLeft := false
			bounceRight := false

			if numberLocation[0] > 0 {

				bounceLeft = true
				adjacentFields = append(adjacentFields, string(rune(line[numberLocation[0]-1])))
			}

			if numberLocation[1] < limitX {

				bounceRight = true
				adjacentFields = append(adjacentFields, string(rune(line[numberLocation[1]])))
			}

			if indexY > 0 {

				if bounceLeft {

					adjacentFields = append(adjacentFields, string(rune(board[indexY-1][numberLocation[0]-1])))
				}

				if bounceRight {

					adjacentFields = append(adjacentFields, string(rune(board[indexY-1][numberLocation[1]])))
				}

				for i := numberLocation[0]; i <= numberLocation[1]; i++ {

					adjacentFields = append(adjacentFields, string(rune(board[indexY-1][i])))
				}
			}

			if indexY < limitY {

				if bounceLeft {

					adjacentFields = append(adjacentFields, string(rune(board[indexY+1][numberLocation[0]-1])))
				}

				if bounceRight {

					adjacentFields = append(adjacentFields, string(rune(board[indexY+1][numberLocation[1]])))
				}

				for i := numberLocation[0]; i < numberLocation[1]; i++ {

					adjacentFields = append(adjacentFields, string(rune(board[indexY+1][i])))
				}
			}

			if containsSymbol(adjacentFields) {

				partNumbers = append(partNumbers, getNumberAt(board, indexY, numberLocation))
			}
		}
	}

	return partNumbers
}

func containsSymbol(fields []string) bool {

	re := regexp.MustCompile(`\d|\.|\n|\t|\r`)

	for _, field := range fields {

		if !re.MatchString(fmt.Sprintf("%s", field)) {

			return true
		}
	}

	return false
}

func getNumberAt(board []string, locY int, locsX []int) int {

	numberString := board[locY][locsX[0]:locsX[1]]

	number, err := strconv.ParseInt(numberString, 10, 32)

	if err != nil {

		panic(err)
	}

	return int(number)
}

func getGearRatios(board []string) []int64 {

	limitY := len(board) - 1
	limitX := len(board[0]) - 1

	var gearRatios []int64

	re := regexp.MustCompile(`\*`)

	for indexY, line := range board {

		gearLocations := re.FindAllStringIndex(line, -1)

		for _, gearLocation := range gearLocations {

			var lines []string

			bounceLeft := false
			bounceRight := false

			if gearLocation[0] > 0 {

				bounceLeft = true

				currentLeft := gearLocation[0] - 1

				for {
					if (currentLeft <= 0) || (!isNumber(rune(line[currentLeft]))) {

						break
					}
					currentLeft--
				}

				lines = append(lines, line[currentLeft:gearLocation[0]])
			}

			if gearLocation[0] < limitX {

				bounceRight = true

				currentRight := gearLocation[0] + 1

				for {
					if (currentRight >= limitX) || (!isNumber(rune(line[currentRight]))) {

						break
					}

					currentRight++
				}

				lines = append(lines, line[gearLocation[0]:currentRight])
			}

			if indexY > 0 {

				currentIndexY := indexY - 1

				currentLeft := gearLocation[0]

				if bounceLeft {

					currentLeft = gearLocation[0] - 1

					for {

						if (currentLeft <= 0) || (!isNumber(rune(board[currentIndexY][currentLeft]))) {

							break
						}

						currentLeft--
					}
				}

				currentRight := gearLocation[0]

				if bounceRight {

					currentRight = gearLocation[0] + 1

					for {

						if (currentRight >= limitX) || (!isNumber(rune(board[currentIndexY][currentRight]))) {

							break
						}
						currentRight++
					}
				}

				lines = append(lines, board[indexY-1][currentLeft:currentRight])
			}

			if indexY < limitY {

				currentIndexY := indexY + 1

				currentLeft := gearLocation[0]

				if bounceLeft {

					currentLeft = gearLocation[0] - 1

					for {

						if (currentLeft <= 0) || (!isNumber(rune(board[currentIndexY][currentLeft]))) {

							break
						}

						currentLeft--
					}
				}

				currentRight := gearLocation[0]

				if bounceRight {

					currentRight = gearLocation[0] + 1

					for {

						if (currentRight >= limitX) || (!isNumber(rune(board[currentIndexY][currentRight]))) {

							break
						}
						currentRight++
					}
				}

				lines = append(lines, board[currentIndexY][currentLeft:currentRight])
			}

			joinedLines := strings.Join(lines, "_")

			reNumbers := regexp.MustCompile(`\d+`)

			matches := reNumbers.FindAllString(joinedLines, -1)

			if len(matches) == 2 {

				numberOne, err1 := strconv.ParseInt(matches[0], 10, 32)

				if err1 != nil {

					panic(err1)
				}

				numberTwo, err2 := strconv.ParseInt(matches[1], 10, 32)

				if err2 != nil {

					panic(err2)
				}

				gearRatios = append(gearRatios, numberOne*numberTwo)
			}
		}
	}

	return gearRatios
}

func isNumber(character rune) bool {

	re := regexp.MustCompile(`\d`)

	characterString := string(character)

	return re.MatchString(characterString)
}

func Part1(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	partNumbers := getPartNumbers(lines)

	var sum int64 = 0

	for _, number := range partNumbers {

		sum += int64(number)
	}

	return strconv.FormatInt(sum, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	gearRatios := getGearRatios(lines)

	var sum int64 = 0

	for _, number := range gearRatios {

		sum += number
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

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
