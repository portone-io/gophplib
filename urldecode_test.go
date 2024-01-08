package gophplib

import (
	"fmt"
	"testing"
)

func ExampleUrldecode() {
	fmt.Println(Urldecode("my=apples&are=green+and%20red%F0%9F%8D%8E"))
	fmt.Println(Urldecode("foo.php?myvar=%BA"))
	// Output:
	// my=apples&are=green and red🍎
	// foo.php?myvar=�
}

// Test cases for Urldecode. These tests were created using the following test
// cases in PHP as inspiration.
//
// Reference:
//   - https://www.php.net/manual/en/function.urldecode.php
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/url/urldecode_variation_001.phpt
func TestUrldecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "emptyString",
			input:    "",
			expected: "",
		},
		{
			name:     "PlusSymbolToSpace",
			input:    "my=apples&are=green+and+red",
			expected: "my=apples&are=green and red",
		},
		{
			name:     "Space",
			input:    "hello%20world",
			expected: "hello world",
		},
		{
			name:     "ExclamationMark",
			input:    "hello%21world",
			expected: "hello!world",
		},
		{
			name:     "NumberSign",
			input:    "number%23",
			expected: "number#",
		},
		{
			name:     "DollarSign",
			input:    "amount%24",
			expected: "amount$",
		},
		{
			name:     "PercentSignEncoded",
			input:    "discount%25",
			expected: "discount%",
		},
		{
			name:     "Ampersand",
			input:    "this%26that",
			expected: "this&that",
		},
		{
			name:     "PlusAsPlus",
			input:    "plus%2Bsign",
			expected: "plus+sign",
		},
		{
			name:     "Comma",
			input:    "comma%2Cseparated",
			expected: "comma,separated",
		},
		{
			name:     "Slash",
			input:    "forward%2Fslash",
			expected: "forward/slash",
		},
		{
			name:     "Colon",
			input:    "time%3A12%3A00",
			expected: "time:12:00",
		},
		{
			name:     "SemiColon",
			input:    "semicolon%3B",
			expected: "semicolon;",
		},
		{
			name:     "Equals",
			input:    "equal%3Dsign",
			expected: "equal=sign",
		},
		{
			name:     "QuestionMark",
			input:    "query%3Fparam",
			expected: "query?param",
		},
		{
			name:     "AtSign",
			input:    "email%40example.com",
			expected: "email@example.com",
		},
		{
			name:     "OpenBracket",
			input:    "open%5Bbracket",
			expected: "open[bracket",
		},
		{
			name:     "CloseBracket",
			input:    "close%5Dbracket",
			expected: "close]bracket",
		},
		{
			name:     "Caret",
			input:    "up%5Ehigh",
			expected: "up^high",
		},
		{
			name:     "MagicQuote",
			input:    "script.php?sterm=%2527",
			expected: "script.php?sterm=%27",
		},
		{
			name:     "Singlequote",
			input:    "single%27quote",
			expected: "single'quote",
		},
		{
			name:     "BypassEncoded",
			input:    "by%0pass",
			expected: "by%0pass",
		},
		{
			name:     "InvalidPercentSign1",
			input:    "red%",
			expected: "red%",
		},
		{
			name:     "InvalidPercentSign2",
			input:    "red%%%blue",
			expected: "red%%%blue",
		},
		{
			name:     "InvalidPercentSign3",
			input:    "red%0blue",
			expected: "red\x0blue",
		},
		{
			name:     "InvalidPercentSign4",
			input:    "red%00green",
			expected: "red\000green",
		},
		{
			name:     "UnicodeEmoji",
			input:    "my=apples&are=green%20and%20red%F0%9F%8D%8E",
			expected: "my=apples&are=green and red🍎",
		},
		{
			name:     "DecodeSingleByteInvalidUTF8",
			input:    "foo.php?myvar=%BA",
			expected: "foo.php?myvar=�",
		},
		{
			name:     "DecodeValidTwoByteUTF8Character",
			input:    "foo.php?myvar=%C2%BA",
			expected: "foo.php?myvar=º",
		},
		{
			name:     "WrongNullChars",
			input:    "yo %0%0%0 lo",
			expected: "yo %0%0%0 lo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Urldecode(tc.input)
			if result != tc.expected {
				t.Errorf(`
expected  %#v
actual    %#v`, tc.expected, result)
			}
		})
	}
}

func TestIsxdigit(t *testing.T) {
	cases := []struct {
		byte
		bool
	}{
		{'0', true},
		{'1', true},
		{'2', true},
		{'3', true},
		{'4', true},
		{'5', true},
		{'6', true},
		{'7', true},
		{'8', true},
		{'9', true},
		{'a', true},
		{'b', true},
		{'c', true},
		{'d', true},
		{'e', true},
		{'f', true},
		{'A', true},
		{'B', true},
		{'C', true},
		{'D', true},
		{'E', true},
		{'F', true},
		{'g', false},
		{'h', false},
		{'+', false},
		{' ', false},
		{'!', false},
		{'z', false},
	}

	for _, tc := range cases {
		t.Run(string(tc.byte), func(t *testing.T) {
			if isxdigit(tc.byte) != tc.bool {
				t.Errorf("isxdigit(%q) != %v", tc.byte, tc.bool)
			}
		})
	}
}

func TestHtoi(t *testing.T) {
	cases := []struct {
		string
		byte
	}{
		{"00", 0x00},
		{"AB", 0xAB},
		{"ab", 0xAB},
		{"FF", 0xFF},
		{"ff", 0xFF},
		{"09", 0x09},
		{"0a", 0x0A},
		{"C0", 0xC0},
	}

	for _, tc := range cases {
		t.Run(tc.string, func(t *testing.T) {
			if htoi(tc.string[0], tc.string[1]) != tc.byte {
				t.Errorf("htoi(%v) != %v", tc.string, tc.byte)
			}
		})
	}
}
