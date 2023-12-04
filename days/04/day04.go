package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const DAY = "04"

type CardData struct {
	card        Card
	unprocessed int
	processed   int
}

type Card struct {
	id             int
	winningNumbers []int
	ownNumbers     []int
}

func (b Card) String() string {

	return fmt.Sprintf("Card(id=%d, winning=%#v, own=%#v)", b.id, b.winningNumbers, b.ownNumbers)
}

func parseCard(line string) Card {

	numberRegex := regexp.MustCompile(`\d+`)

	cardString := strings.Split(line, ":")[0]

	numbersString := strings.Split(line, ":")[1]

	winningNumbersString := strings.Split(numbersString, "|")[0]
	ownNumbersString := strings.Split(numbersString, "|")[1]

	var card Card

	card.id = int(stringToNumber(numberRegex.FindString(cardString)))

	allWinningStrings := numberRegex.FindAllString(winningNumbersString, -1)

	for _, winningNumberString := range allWinningStrings {

		card.winningNumbers = append(card.winningNumbers, int(stringToNumber(winningNumberString)))
	}

	allOwnNumberStrings := numberRegex.FindAllString(ownNumbersString, -1)

	for _, ownNumberString := range allOwnNumberStrings {

		card.ownNumbers = append(card.ownNumbers, int(stringToNumber(ownNumberString)))
	}

	return card
}

func stringToNumber(s string) int64 {

	number, err := strconv.ParseInt(s, 10, 32)

	if err != nil {

		panic(err)
	}

	return number
}

func getMatchingNumbers(card Card) []int {

	var matchingNumbers []int

	for _, ownNumber := range card.ownNumbers {

		if slices.Contains(card.winningNumbers, ownNumber) {

			matchingNumbers = append(matchingNumbers, ownNumber)
		}
	}

	return matchingNumbers
}

func getWorth(card Card) int64 {

	var worth int64 = 0

	for i := 0; i < len(getMatchingNumbers(card)); i++ {

		if worth == 0 {

			worth = 1
		} else {

			worth *= 2
		}
	}

	return worth
}

func Part1(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	var totalWorth int64 = 0

	for _, line := range lines {

		card := parseCard(line)
		totalWorth += getWorth(card)
	}

	return strconv.FormatInt(totalWorth, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	lines := strings.Split(content, "\n")

	cardMap := make(map[int]CardData)

	processingQueue := make([]int, 0)

	for _, line := range lines {

		card := parseCard(line)
		cardMap[card.id] = CardData{card: card, unprocessed: 1, processed: 0}
		processingQueue = append(processingQueue, card.id)
	}

	for {

		if len(processingQueue) == 0 {

			break
		}

		cardId := processingQueue[0]
		processingQueue = processingQueue[1:]

		card := cardMap[cardId]
		matches := len(getMatchingNumbers(card.card))

		for i := cardId + 1; i <= cardId+matches; i++ {

			matchCard := cardMap[i]
			matchCard.unprocessed += card.unprocessed
			if !slices.Contains(processingQueue, i) {
				processingQueue = append(processingQueue, i)
			}
			cardMap[i] = matchCard
		}

		card.processed += card.unprocessed
		card.unprocessed = 0

		cardMap[cardId] = card
	}

	var totalSum int64 = 0

	for _, cardData := range cardMap {

		totalSum += int64(cardData.processed)
	}

	return strconv.FormatInt(totalSum, 10)
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
