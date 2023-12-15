package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const DAY = "15"

type Boxes struct {
	numbersToBox map[int]Box
}

type Box struct {
	// contains a mapping of at least the labels in the box to their focal lengths
	labelsToFocalLengths map[string]int
	// contains the sorted labels in the box separated by comma
	sortedLabelsCommaSeperated string
}

func sumHash(input []string) int {

	sum := 0

	for _, in := range input {

		sum += hash(in)
	}

	return sum
}

func hash(input string) int {

	currentValue := 0

	for _, character := range input {

		currentValue += int(character)
		currentValue *= 17
		currentValue = int(math.Mod(float64(currentValue), 256))
	}

	return currentValue
}

func parseInput(input string) []string {

	return strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "\r", ""), "\n", ""), " ", ""), ",")
}

func applyLens(lens string, boxes *Boxes) {

	// label is everything before "-" or "="
	label := strings.Split(strings.Split(lens, "=")[0], "-")[0]

	// command is the first character after the label
	command := lens[len(label) : len(label)+1]

	hash := hash(label)

	box, boxExists := boxes.numbersToBox[hash]

	if command == "-" {

		if boxExists {

			box.sortedLabelsCommaSeperated = strings.Replace(box.sortedLabelsCommaSeperated, label+",", "", 1)
		}

	} else if command == "=" {

		focalLength := stringToNumber(lens[len(lens)-1:])

		if !boxExists {

			box = Box{make(map[string]int), ""}
		}

		if box.labelsToFocalLengths == nil {

			box.labelsToFocalLengths = make(map[string]int)
		}

		lensExists := strings.Contains(box.sortedLabelsCommaSeperated, label)

		box.labelsToFocalLengths[label] = int(focalLength)

		if !lensExists {

			box.sortedLabelsCommaSeperated += label + ","
		}
	}

	boxes.numbersToBox[hash] = box
}

func getFocussingPower(boxes Boxes) int {

	sum := 0

	for boxNumber, box := range boxes.numbersToBox {

		sum += getFocussingPowerBox(box, boxNumber)
	}

	return sum
}

func getFocussingPowerBox(box Box, boxNumber int) int {

	sum := 0

	if len(box.sortedLabelsCommaSeperated) > 0 {

		for lensIndex, lens := range strings.Split(box.sortedLabelsCommaSeperated[:len(box.sortedLabelsCommaSeperated)-1], ",") {

			focalLength, lensExists := box.labelsToFocalLengths[lens]

			if lensExists {

				sum += (1 + boxNumber) * (1 + lensIndex) * focalLength
			}
		}
	}

	return sum
}

func Part1(input string) string {

	content := GetContent(input)

	result := sumHash(parseInput(content))

	return strconv.Itoa(result)
}

func Part2(input string) string {

	content := GetContent(input)

	sequence := parseInput(content)

	boxes := Boxes{numbersToBox: make(map[int]Box)}

	for _, lens := range sequence {

		applyLens(lens, &boxes)
	}

	result := getFocussingPower(boxes)

	return strconv.Itoa(result)
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
