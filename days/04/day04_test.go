package main

import (
	"fmt"
	"strconv"
	"testing"
)

var P1_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P1_OUT_TEST = []string{"13"}

var P2_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P2_OUT_TEST = []string{"30"}

func TestParseCard(t *testing.T) {

	testIn := "Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19"
	testOut := Card{id: 2, winningNumbers: []int{13, 32, 20, 16, 61}, ownNumbers: []int{61, 30, 68, 82, 17, 32, 24, 19}}

	testCard := parseCard(testIn)

	assert("parseCard", testIn, fmt.Sprintf("%s", testOut), fmt.Sprintf("%s", testCard), t)
}

func TestWorth(t *testing.T) {

	testIn := parseCard("Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53")
	var testOut int64 = 8

	result := getWorth(testIn)

	assert("getWorth", fmt.Sprintf("%s", testIn), strconv.FormatInt(testOut, 10), strconv.FormatInt(result, 10), t)
}

func TestPart1(t *testing.T) {

	for index, element := range P1_IN_TEST {

		expected := P1_OUT_TEST[index]
		received := Part1(element)
		assert("Part1", element, expected, received, t)
	}
}

func TestPart2(t *testing.T) {

	for index, element := range P2_IN_TEST {

		expected := P2_OUT_TEST[index]
		received := Part2(element)
		assert("Part2", element, expected, received, t)
	}
}

func assert(method string, input string, expected string, received string, t *testing.T) {

	if expected != received {

		t.Errorf("%s(%s) expected '%s' but received '%s'", method, input, expected, received)
	}
}
