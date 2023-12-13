package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "13"

type Data struct {
	rows    []string
	columns []string
}

func getNumberOfLinesVertical(data Data, smudges int) int64 {

	sum := int64(0)

	for reflectionLineIndex := 0; reflectionLineIndex < len(data.columns)-1; reflectionLineIndex++ {

		if checkHorizontalReflection(data, reflectionLineIndex, true, smudges) {

			sum += int64(reflectionLineIndex + 1)
		}
	}

	return sum
}

func getNumberOfLinesHorizontal(data Data, smudges int) int64 {

	sum := int64(0)

	for reflectionLineIndex := 0; reflectionLineIndex < len(data.rows)-1; reflectionLineIndex++ {

		if checkHorizontalReflection(data, reflectionLineIndex, false, smudges) {

			sum += int64(reflectionLineIndex + 1)
		}
	}

	return sum
}

func checkHorizontalReflection(data Data, reflectionLineIndex int, isColumn bool, smudges int) bool {

	items := data.rows

	if isColumn {

		items = data.columns
	}

	topIndex := reflectionLineIndex
	bottomIndex := reflectionLineIndex + 1

	foundSmudges := 0

	for {
		if topIndex < 0 || bottomIndex > len(items)-1 {

			break
		}

		if items[topIndex] != items[bottomIndex] {

			for charIndex := 0; charIndex < len(items[topIndex]); charIndex++ {

				if items[topIndex][charIndex] != items[bottomIndex][charIndex] {

					foundSmudges++
				}
			}

			if foundSmudges > smudges {

				return false
			}
		}

		topIndex--
		bottomIndex++
	}

	if foundSmudges == smudges {

		return true
	} else {

		return false
	}
}

func parseData(input string) Data {

	var data Data

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	for _, line := range lines {

		data.rows = append(data.rows, line)
	}

	for x := 0; x < len(lines[0]); x++ {

		line := ""
		for y := 0; y < len(lines); y++ {

			line = line + string(lines[y][x])
		}

		data.columns = append(data.columns, line)
	}

	return data
}

func Part1(input string) string {

	content := GetContent(input)

	chunks := strings.Split(strings.ReplaceAll(content, "\r", ""), "\n\n")

	sum := int64(0)

	for _, chunk := range chunks {

		data := parseData(chunk)

		sum += getNumberOfLinesVertical(data, 0)
		sum += getNumberOfLinesHorizontal(data, 0) * 100
	}

	return strconv.FormatInt(sum, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	chunks := strings.Split(strings.ReplaceAll(content, "\r", ""), "\n\n")

	sum := int64(0)

	for _, chunk := range chunks {

		data := parseData(chunk)

		sum += getNumberOfLinesVertical(data, 1)
		sum += getNumberOfLinesHorizontal(data, 1) * 100
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
