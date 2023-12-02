package main

import (
	"fmt"
	"testing"
)

var P1_IN_TEST = [1]string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P1_OUT_TEST = [1]string{"8"}

var P2_IN_TEST = [1]string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P2_OUT_TEST = [1]string{"2286"}

func TestParseColor(t *testing.T) {

	testIn := [4]string{"red", "green", "blue", "yellow"}
	testOut := [4]Color{Red, Green, Blue, NotAColor}

	for index, test := range testIn {

		assert("parseColor", test, fmt.Sprintf("%d", testOut[index]), fmt.Sprintf("%d", parseColor(test)), t)
	}
}

func TestParseBallCount(t *testing.T) {

	testIn := "  8 green "
	testOut := BallCount{count: 8, color: Green}

	ballCountReceived := parseBallCount(testIn)

	assert("parseBallCount", testIn, fmt.Sprintf("%s", testOut), fmt.Sprintf("%s", ballCountReceived), t)
}

func TestParseDraw(t *testing.T) {

	testIn := "  1 red, 2 green, 6 blue "
	testOut := Draw{ballCounts: []BallCount{{count: 1, color: Red}, {count: 2, color: Green}, {count: 6, color: Blue}}}

	drawReceived := parseDraw(testIn)

	assert("parseDraw", testIn, fmt.Sprintf("%s", drawReceived), fmt.Sprintf("%s", testOut), t)
}

func TestParseGame(t *testing.T) {

	testIn := "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red"
	testOut := Game{
		id: 3,
		draws: []Draw{
			{
				ballCounts: []BallCount{
					{count: 8, color: Green},
					{count: 6, color: Blue},
					{count: 20, color: Red}},
			}, {
				ballCounts: []BallCount{
					{count: 5, color: Blue},
					{count: 4, color: Red},
					{count: 13, color: Green}},
			}, {
				ballCounts: []BallCount{
					{count: 5, color: Green},
					{count: 1, color: Red},
				},
			}}}

	gameReceived := parseGame(testIn)

	assert("parseGame", testIn, fmt.Sprintf("%s", testOut), fmt.Sprintf("%s", gameReceived), t)
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
