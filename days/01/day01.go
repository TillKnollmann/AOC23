package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DAY = "01"

func getLineCalibration(line string) int {

	characters := strings.Split(line, "")

	first := -1
	last := -1

	for _, character := range characters {

		intValue, err := strconv.ParseInt(character, 10, 8)

		if err == nil {

			if first == -1 {

				first = int(intValue)
			} else {

				last = int(intValue)
			}
		}
	}

	if last == -1 {

		last = first
	}

	return combine(first, last)
}

func convertToInt(match string) int {

	value, err := strconv.Atoi(match)
	if err != nil {
		switch match {
		case "one":
			value = 1
		case "two":
			value = 2
		case "three":
			value = 3
		case "four":
			value = 4
		case "five":
			value = 5
		case "six":
			value = 6
		case "seven":
			value = 7
		case "eight":
			value = 8
		case "nine":
			value = 9
		default:
		}
	}
	return value
}

func convertToIntReversed(match string) int {

	value, err := strconv.Atoi(match)
	if err != nil {
		switch match {
		case "eno":
			value = 1
		case "owt":
			value = 2
		case "eerht":
			value = 3
		case "ruof":
			value = 4
		case "evif":
			value = 5
		case "xis":
			value = 6
		case "neves":
			value = 7
		case "thgie":
			value = 8
		case "enin":
			value = 9
		default:
		}
	}
	return value
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func combine(val1 int, val2 int) int {

	resString := fmt.Sprintf("%d%d", val1, val2)

	resInt, _ := strconv.ParseInt(resString, 10, 8)

	return int(resInt)
}

func getLineCalibrationPart2(line string) int {

	re := regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine`)
	first := convertToInt(re.FindString(line))

	reReversed := regexp.MustCompile(`\d|enin|thgie|neves|xis|evif|ruof|eerht|owt|eno`)
	lineReversed := Reverse(line)

	last := convertToIntReversed(reReversed.FindString(lineReversed))

	return combine(first, last)
}

func Part1(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	sum := 0

	for _, element := range lines {

		sum += getLineCalibration(element)
	}

	return fmt.Sprint(sum)
}

func Part2(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	sum := 0

	for _, element := range lines {

		sum += getLineCalibrationPart2(element)
	}

	return fmt.Sprint(sum)
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
