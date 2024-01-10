package gophplib

import (
	"fmt"
	"reflect"
	"testing"
)

var unsetVariable string

func ExampleOrd() {
	// Plain char
	fmt.Println(Ord("a"))

	// Plain string
	fmt.Println(Ord("Hello"))

	// Special character
	fmt.Println(Ord("!"))

	// New line character
	fmt.Println(Ord("\n"))

	// Hexadecimal representation of newline character
	fmt.Println(Ord("\x0A"))

	// Character '0'
	fmt.Println(Ord("0"))

	// Emoji
	fmt.Println(Ord("🐘"))

	// 한글
	fmt.Println(Ord("안녕하세요"))

	// Output:
	// 97 <nil>
	// 72 <nil>
	// 33 <nil>
	// 10 <nil>
	// 10 <nil>
	// 48 <nil>
	// 240 <nil>
	// 236 <nil>
}

func ExampleOrd_variation() {
	// Array
	fmt.Println(Ord([]int{1}))

	// Nil
	fmt.Println(Ord(nil))

	// Empty string
	fmt.Println(Ord(""))

	// Unset variable
	fmt.Println(Ord(unsetVariable))

	// Output:
	// 0 unsupported type : []int
	// 0 <nil>
	// 0 <nil>
	// 0 <nil>
}

func TestOrd(t *testing.T) {
	testCases := []struct {
		testName string
		input    any
		expected byte
	}{
		{
			testName: "Plain char",
			input:    "a",
			expected: 97,
		},
		{
			testName: "Plain string",
			input:    "Hello",
			expected: 72,
		},
		{
			testName: "Special character",
			input:    "!",
			expected: 33,
		},
		{
			testName: "New line character",
			input:    "\n",
			expected: 10,
		},
		{
			testName: "Hexadecimal representation of newline character",
			input:    "\x0A",
			expected: 10,
		},
		{
			testName: "Hexadecimal representation of newline character",
			input:    "\x0A",
			expected: 10,
		},
		{
			testName: "Character '0'",
			input:    "0",
			expected: 48,
		},
		{
			testName: "Array",
			input:    []int{1},
			expected: 0,
		},
		{
			testName: "Nil",
			input:    nil,
			expected: 0,
		},
		{
			testName: "Empty string",
			input:    "",
			expected: 0,
		},
		{
			testName: "Unset variable",
			input:    unsetVariable,
			expected: 0,
		},
		{
			testName: "Emoji",
			input:    "🐘",
			expected: 240,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := Ord(tc.input)
			if err != nil {
				expectedErr := fmt.Errorf("unsupported type : %s", reflect.TypeOf(tc.input))
				if err.Error() != expectedErr.Error() {
					t.Errorf("%s: expected error : %s, got %s", tc.testName, expectedErr, err)
				}
			} else {
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("%s: expected %v, got %v", tc.testName, tc.expected, result)
				}
			}
		})
	}
}
