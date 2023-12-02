package main

import (
	"fmt"
	"testing"
)

var P1_IN_TEST = [1]string{"test/02/in01.txt"}
var P1_OUT_TEST = [1]string{"8"}

var P2_IN_TEST = [1]string{"test/02/in01.txt"}
var P2_OUT_TEST = [1]string{"2286"}

func TestParseColor(t *testing.T) {

	testIn := [4]string{"red", "green", "blue", "yellow"}
	testOut := [4]Color{Red, Green, Blue, NotAColor}

	for index, test := range testIn {

		if testOut[index] != parseColor(test) {

			t.Errorf("parseColor(%s) expected '%d' received '%d'", test, testOut[index], parseColor(test))
		}
	}
}

func TestParseBallCount(t *testing.T) {

	testIn := "  8 green "
	testOut := BallCount{count: 8, color: Green}

	ballCountReceived := parseBallCount(testIn)

	if ballCountReceived.count != testOut.count || ballCountReceived.color != testOut.color {

		t.Errorf("parseBallCount(%s) expected '%s' received '%s'", testIn, testOut, ballCountReceived)
	}
}

func TestParseDraw(t *testing.T) {

	testIn := "  1 red, 2 green, 6 blue "
	testOut := Draw{ballCounts: []BallCount{{count: 1, color: Red}, {count: 2, color: Green}, {count: 6, color: Blue}}}

	drawReceived := parseDraw(testIn)

	if fmt.Sprintf("%s", drawReceived) != fmt.Sprintf("%s", testOut) {

		t.Errorf("parseDraw(%s) expected '%s' received '%s'", testIn, testOut, drawReceived)
	}
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

	if fmt.Sprintf("%s", gameReceived) != fmt.Sprintf("%s", testOut) {

		t.Errorf("parseGame(%s) expected '%s' received '%s'", testIn, testOut, gameReceived)
	}
}

func TestPart1(t *testing.T) {

	for index, element := range P1_IN_TEST {

		expected := P1_OUT_TEST[index]
		received := Part1(element)
		assert(expected, received, t)
	}
}

func TestPart2(t *testing.T) {

	for index, element := range P2_IN_TEST {

		expected := P2_OUT_TEST[index]
		received := Part2(element)
		assert(expected, received, t)
	}
}

func assert(expected string, received string, t *testing.T) {

	if expected != received {

		t.Errorf("Expected '%s' but received '%s'", expected, received)
	}
}
