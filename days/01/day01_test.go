package main

import "testing"

var P1_IN_TEST = [1]string{"test/01/in01.txt"}
var P1_OUT_TEST = [1]string{"142"}

var P2_IN_TEST = [2]string{"test/01/in02.txt", "test/01/in03.txt"}
var P2_OUT_TEST = [2]string{"281", "58"}

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
