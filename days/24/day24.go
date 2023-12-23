package main

import (
	"fmt"
	"log"
	"os"
)

const DAY = "24"

func Part1(input string, maxSteps int) string {

	content := GetContent(input)

	return string(content[0])
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

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY), 64)))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
