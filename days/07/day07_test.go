package main

import (
	"fmt"
	"testing"
)

var P1_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P1_OUT_TEST = []string{"6440"}

var P2_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P2_OUT_TEST = []string{"5905"}

func TestParseHand(t *testing.T) {

	testIn := "32T3K 765"
	testOut := Hand{bid: 765, stringRep: "32T3K", handType: OnePair}

	result := parseHand(testIn, false)

	assert("parseHand", testIn, fmt.Sprintf("%s", testOut), fmt.Sprintf("%s", result), t)
}

func TestParseHandType(t *testing.T) {

	testIn := []string{"AAAAA", "AA8AA", "23332", "TTT98", "23432", "A23A4", "23456"}
	testOut := []HandType{FiveOfAKind, FourOfAKind, FullHouse, ThreeOfAKind, TwoPair, OnePair, HighCard}

	for index := range testIn {

		assert("parseHandType", testIn[index], fmt.Sprintf("%d", testOut[index]), fmt.Sprintf("%d", parseHandType(testIn[index], false)), t)
	}
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
