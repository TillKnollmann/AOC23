package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const DAY = "07"

type Hand struct {
	bid       int64
	handType  HandType
	stringRep string
}

func (b Hand) String() string {

	return fmt.Sprintf("Hand(bid=%d, handType=%d, stringRep=%s)", b.bid, b.handType, b.stringRep)
}

var CardValues = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var AllCardTypesExceptJoker = []string{
	"2", "3", "4", "5", "6", "7", "8", "9", "T", "Q", "K", "A",
}

type HandType int

const (
	WTFisThisCard HandType = iota + 1
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func parseAllHands(input string, allowJoker bool) []Hand {

	var hands []Hand

	for _, line := range strings.Split(input, "\n") {

		hands = append(hands, parseHand(line, allowJoker))
	}

	return hands
}

func parseHand(input string, allowJoker bool) Hand {

	var hand Hand

	numberRe := regexp.MustCompile(`\d+`)

	hand.bid = stringToNumber(numberRe.FindString(strings.Split(input, " ")[1]))
	hand.stringRep = strings.Split(input, " ")[0]
	hand.handType = parseHandType(hand.stringRep, allowJoker)

	return hand
}

func parseHandType(input string, allowJoker bool) HandType {

	handString := input
	charMap := getCharOccurrences(handString)
	charValues := getValues(charMap)

	if !allowJoker || !strings.Contains(handString, "J") {

		if isFiveOfAKind(charMap) {
			return FiveOfAKind
		} else if isFourOfAKind(charValues) {
			return FourOfAKind
		} else if isFullHouse(charValues) {
			return FullHouse
		} else if isThreeOfAKind(charValues) {
			return ThreeOfAKind
		} else if isTwoPair(charValues) {
			return TwoPair
		} else if isOnePair(charValues) {
			return OnePair
		} else if isHighCard(charValues) {
			return HighCard
		} else {
			return WTFisThisCard
		}
	}

	if charMap['J'] >= 4 || len(charValues) == 2 {
		return FiveOfAKind
	} else if charMap['J'] == 3 {
		return FourOfAKind
	} else if charMap['J'] == 2 {

		var highestHandType HandType
		for _, firstJoker := range AllCardTypesExceptJoker {

			for _, secondJoker := range AllCardTypesExceptJoker {
				highestHandType = max(highestHandType, parseHandType(strings.Replace(strings.Replace(handString, "J", firstJoker, 1), "J", secondJoker, 1), false))
			}
		}

		return highestHandType

	} else if charMap['J'] == 1 {

		var highestHandType HandType
		for _, joker := range AllCardTypesExceptJoker {

			highestHandType = max(highestHandType, parseHandType(strings.Replace(handString, "J", joker, -1), false))
		}

		return highestHandType
	}

	return WTFisThisCard
}

func isHighCard(values []int) bool {

	return len(values) == 5
}

func isOnePair(values []int) bool {

	if len(values) == 4 {

		return true
	}
	return false
}

func isTwoPair(values []int) bool {

	if len(values) == 3 {

		if values[0] == 2 && values[1] == 2 {

			return true
		}
	}
	return false
}

func isThreeOfAKind(values []int) bool {

	if len(values) == 3 {

		if slices.Contains(values, 3) && slices.Contains(values, 1) {

			return true
		}
	}
	return false
}

func isFullHouse(values []int) bool {

	if len(values) == 2 {

		if slices.Contains(values, 3) {

			return true
		}
	}
	return false
}

func isFourOfAKind(charValues []int) bool {

	if len(charValues) == 2 {

		if slices.Contains(charValues, 4) {

			return true
		}
	}
	return false
}

func isFiveOfAKind(charMap map[rune]int) bool {

	return len(charMap) == 1
}

func getValues(m map[rune]int) []int {

	list := make([]int, 0, len(m))

	for _, value := range m {
		list = append(list, value)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})

	return list
}

func getCharOccurrences(input string) map[rune]int {

	result := map[rune]int{}
	for _, char := range input {
		result[char] += 1
	}

	return result
}

func sortHands(hands []Hand, allowJoker bool) []Hand {

	var newHands []Hand = make([]Hand, len(hands))

	copy(newHands, hands)

	sort.Slice(newHands, func(i, j int) bool {
		return isFirstSmallerThanSecond(newHands[i], newHands[j], allowJoker)
	})

	return newHands
}

func isFirstSmallerThanSecond(first Hand, second Hand, allowJoker bool) bool {

	if first.handType < second.handType {

		return true
	} else if first.handType > second.handType {

		return false
	}

	if allowJoker {
		CardValues['J'] = 1
	}

	for index := range first.stringRep {

		if CardValues[rune(first.stringRep[index])] < CardValues[rune(second.stringRep[index])] {

			return true
		} else if CardValues[rune(first.stringRep[index])] > CardValues[rune(second.stringRep[index])] {

			return false
		}
	}

	panic(fmt.Sprintf("Cards %s and %s are identical", first.stringRep, second.stringRep))
}

func Part1(input string) string {

	content := GetContent(input)

	hands := parseAllHands(content, false)
	hands = sortHands(hands, false)

	var result int64

	for index, hand := range hands {

		rank := index + 1
		result += int64(rank) * hand.bid
	}

	return strconv.FormatInt(result, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	hands := parseAllHands(content, true)
	hands = sortHands(hands, true)

	var result int64

	for index, hand := range hands {

		rank := index + 1
		result += int64(rank) * hand.bid
	}

	return strconv.FormatInt(result, 10)
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
