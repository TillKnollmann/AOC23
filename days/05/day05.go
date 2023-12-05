package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const DAY = "05"

type GardenObject struct {
	name   string
	number int64
}

type Data struct {
	seeds           []int64
	sourceToMapping map[string]Mapping
	minNumber       int64
	maxNumber       int64
}

func (b Data) String() string {

	return fmt.Sprintf("Data(seeds=%#v, sourceToMapping=%#v)", b.seeds, b.sourceToMapping)
}

type Mapping struct {
	source       string
	target       string
	mappingLines []MappingLine
}

type MappingLine struct {
	destRangeStart   int64
	sourceRangeStart int64
	length           int64
}

func (b Mapping) String() string {

	return fmt.Sprintf("Mapping(source=%s, targert=%s, sourceToTarget=%#v)", b.source, b.target, b.mappingLines)
}

func (b MappingLine) String() string {

	return fmt.Sprintf("MappingLine(destRangeStart=%d, sourceRangeStart=%d, length=%d)", b.destRangeStart, b.sourceRangeStart, b.length)
}

func parseData(input string) Data {

	var data Data

	data.sourceToMapping = make(map[string]Mapping)

	data.minNumber = math.MinInt64
	data.maxNumber = math.MaxInt64

	chunks := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")

	numberRe := regexp.MustCompile(`\d+`)

	seedStrings := numberRe.FindAllString(chunks[0], -1)

	for _, seedString := range seedStrings {

		data.seeds = append(data.seeds, stringToNumber(seedString))
	}

	for _, chunk := range chunks[1:] {

		mapping := parseMapping(chunk)
		data.sourceToMapping[mapping.source] = mapping

		for _, mappingLine := range mapping.mappingLines {

			data.minNumber = min(data.minNumber, mappingLine.sourceRangeStart)
			data.maxNumber = max(data.maxNumber, mappingLine.sourceRangeStart+mappingLine.length)
		}
	}

	return data
}

func parseMapping(input string) Mapping {

	var mapping Mapping

	lines := strings.Split(input, "\n")

	sourceTargetStrings := strings.Split(strings.ReplaceAll(lines[0], " map:", ""), "-to-")
	mapping.source = sourceTargetStrings[0]
	mapping.target = sourceTargetStrings[1]

	for _, line := range lines[1:] {

		mapping.mappingLines = append(mapping.mappingLines, parseMappingLine(line))
	}

	return mapping
}

func parseMappingLine(line string) MappingLine {

	var mappingLine MappingLine

	numberRe := regexp.MustCompile(`\d+`)

	numberStrings := numberRe.FindAllString(line, 3)

	mappingLine.destRangeStart = stringToNumber(numberStrings[0])
	mappingLine.sourceRangeStart = stringToNumber(numberStrings[1])
	mappingLine.length = stringToNumber(numberStrings[2])

	return mappingLine
}

func stringToNumber(s string) int64 {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return number
}

func resolve(gardenObject GardenObject, target string, data Data) GardenObject {

	for {
		if gardenObject.name == target {
			break
		}

		gardenObject = resolveStep(gardenObject, data)
	}

	return gardenObject
}

func resolveStep(gardenObject GardenObject, data Data) GardenObject {

	var newGardenObject GardenObject

	newGardenObject.name = data.sourceToMapping[gardenObject.name].target

	found := false

	for _, mappingLine := range data.sourceToMapping[gardenObject.name].mappingLines {

		if found {
			break
		}

		if mappingLine.sourceRangeStart <= gardenObject.number && gardenObject.number < mappingLine.sourceRangeStart+mappingLine.length {

			found = true
			newGardenObject.number = mappingLine.destRangeStart + (gardenObject.number - mappingLine.sourceRangeStart)
		}
	}

	if !found {

		newGardenObject.number = gardenObject.number
	}

	return newGardenObject
}

func Part1(input string) string {

	content := GetContent(input)

	data := parseData(content)

	var minLocationNumber = int64(math.MaxInt64)

	for _, seedNumber := range data.seeds {

		locationNumber := resolve(GardenObject{number: seedNumber, name: "seed"}, "location", data).number
		minLocationNumber = min(minLocationNumber, locationNumber)
	}

	return fmt.Sprintf("%d", minLocationNumber)
}

func Part2(input string) string {

	content := GetContent(input)

	data := parseData(content)

	var minLocationNumber = int64(math.MaxInt64)

	var wg sync.WaitGroup
	wg.Add(len(data.seeds) / 2)

	var lock sync.Mutex

	for index := 0; index < len(data.seeds); index += 2 {

		go func(index int) {

			defer wg.Done()

			for seedNumber := data.seeds[index]; seedNumber < data.seeds[index]+data.seeds[index+1]; seedNumber++ {

				if seedNumber < data.minNumber {

					seedNumber = data.minNumber - 1
				} else if data.maxNumber <= seedNumber {

					seedNumber = data.seeds[index] + data.seeds[index+1] - 1
				} else {

					locationNumber := resolve(GardenObject{number: seedNumber, name: "seed"}, "location", data).number

					if locationNumber < minLocationNumber {
						lock.Lock()
						if locationNumber < minLocationNumber {
							minLocationNumber = min(minLocationNumber, locationNumber)
							fmt.Println(fmt.Sprintf("New minimum: %d", minLocationNumber))
						}
						lock.Unlock()
					}
				}

			}
		}(index)

	}

	wg.Wait()

	return fmt.Sprintf("%d", minLocationNumber)
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
