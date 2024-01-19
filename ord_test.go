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
		any
		byte
	}{
		{"a", 97},
		{"Hello", 72},
		{"!", 33},
		{"\n", 10},
		{"\x0A", 10},
		{"\x0A", 10},
		{"0", 48},
		{[]int{1}, 0},
		{nil, 0},
		{"", 0},
		{unsetVariable, 0},
		{"🐘", 240},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := Ord(tc.any)
			if err != nil {
				byteErr := fmt.Errorf("unsupported type : %s", reflect.TypeOf(tc.any))
				if err.Error() != byteErr.Error() {
					t.Errorf("%s: byte error : %s, got %s", testName, byteErr, err)
				}
			} else {
				if !reflect.DeepEqual(result, tc.byte) {
					t.Errorf("%s: byte %v, got %v", testName, tc.byte, result)
				}
			}
		})
	}
}
