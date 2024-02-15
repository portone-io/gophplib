package gophplib

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleStrlen() {
	// Plain string
	fmt.Println(Strlen("Hello, world"))

	// Special characters
	fmt.Println(Strlen("$@#%^&*!~,.:;?"))

	// Empty string
	fmt.Println(Strlen(""))

	// Empty string with white space
	fmt.Println(Strlen(" "))

	// Nil
	fmt.Println(Strlen(nil))

	// Hexadecimal characters
	fmt.Println(Strlen("\x90\x91\x00\x93\x94\x90\x91\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f"))

	// Tab character
	fmt.Println(Strlen("\t"))

	// '\t' character
	fmt.Println(Strlen(`\t`))

	// Int
	fmt.Println(Strlen(123))

	// Float
	fmt.Println(Strlen(-1.2344))

	// True
	fmt.Println(Strlen(true))

	// False
	fmt.Println(Strlen(false))

	// Byte string
	ac := string([]byte{128, 234, 65, 255, 0}) // chr(128).chr(234).chr(65).chr(255).chr(256)와 동일한 문자
	fmt.Println(Strlen(ac))

	// Multi-byte string
	fmt.Println(Strlen("안녕하세요"))

	// String contain NULL byte
	fmt.Println(Strlen("abc\000def"))

	// Output:
	// 12 <nil>
	// 14 <nil>
	// 0 <nil>
	// 1 <nil>
	// 0 <nil>
	// 18 <nil>
	// 1 <nil>
	// 2 <nil>
	// 3 <nil>
	// 7 <nil>
	// 1 <nil>
	// 0 <nil>
	// 5 <nil>
	// 15 <nil>
	// 7 <nil>
}

func TestStrlen(t *testing.T) {
	testCases := []struct {
		any
		int
	}{
		{"Hello, world", 12},
		{"$@#%^&*!~,.:;?", 14},
		{"", 0},
		{" ", 1},
		{nil, 0},
		{"\x90\x91\x00\x93\x94\x90\x91\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f", 18},
		{"\t", 1},
		{`\t`, 2},
		{123, 3},
		{-1.2344, 7},
		{true, 1},
		{false, 0},
		{string([]byte{128, 234, 65, 255, 0}), 5},
		{"안녕하세요", 15},
		{"abc\000def", 7},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := Strlen(tc.any)
			if err != nil {
				expectedErr := fmt.Errorf("unsupported type : %s", reflect.TypeOf(tc.any))
				if err.Error() != expectedErr.Error() {
					t.Errorf("%s: expected error : %s, got %s", testName, expectedErr, err)
				}
			} else {
				if !reflect.DeepEqual(result, tc.int) {
					t.Errorf("%s: length of input is %d, got %d", testName, tc.int, result)
				}
			}
		})
	}
}

func TestStrlenError(t *testing.T) {
	_, err := Strlen([]int{1, 2, 3})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != "unsupported type : []int" {
		t.Errorf("expected error : unsupported type : []int, got %s", err)
	}
}
