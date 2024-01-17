package gophplib

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
)

type Sample struct{}

type Sample2 struct{}

type CustomType struct {
	value string
}

func (s Sample) toString() string {
	return "sample object"
}

func TestZendParseArgAsString(t *testing.T) {
	file, osErr := os.Open("README.md")
	if osErr != nil {
		panic(osErr)
	}

	testCase := []struct {
		testName string
		input    any
		expected string
	}{
		{
			testName: "Plain string",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			testName: "Empty string",
			input:    "",
			expected: "",
		},
		{
			testName: "StringWithSpecialChars",
			input:    "Line1\nLine2\tTab",
			expected: "Line1\nLine2\tTab",
		},
		{
			testName: "Int",
			input:    123,
			expected: "123",
		},
		{
			testName: "Int64",
			input:    9223372036854775807,
			expected: "9223372036854775807",
		},
		{
			testName: "Negative int",
			input:    -123,
			expected: "-123",
		},
		{
			testName: "Float64",
			input:    123.456,
			expected: "123.456",
		},
		{
			testName: "Float64 exceeds 14 digits",
			input:    123.456789012345678,
			expected: "123.45678901235",
		},
		{
			testName: "Exponent",
			input:    10.1234567e10,
			expected: "101234567000",
		},
		{
			testName: "Special Float - NaN",
			input:    math.NaN(),
			expected: "NAN",
		},
		{
			testName: "Special float - positive infinity",
			input:    math.Inf(1),
			expected: "INF",
		}, {
			testName: "Special float - negative infinity",
			input:    math.Inf(-1),
			expected: "-INF",
		},
		{
			testName: "Exponent",
			input:    10.1234567e10,
			expected: "101234567000",
		},
		{
			testName: "Negative float64",
			input:    -123.456,
			expected: "-123.456",
		},
		{
			testName: "Zero float64",
			input:    0.0,
			expected: "0",
		},
		{
			testName: "True",
			input:    true,
			expected: "1",
		},
		{
			testName: "False",
			input:    false,
			expected: "",
		},
		{
			testName: "Nil",
			input:    nil,
			expected: "",
		},
		{
			testName: "Object with toString function",
			input:    Sample{},
			expected: "sample object",
		},
		{
			testName: "Object without toString function",
			input:    Sample2{},
			expected: "",
		},
		{
			testName: "Int array",
			input:    []int{1, 2, 3},
			expected: "",
		},
		{
			testName: "String array",
			input:    []string{"hello", "world"},
			expected: "",
		},
		{
			testName: "NestedArray",
			input:    []interface{}{[]interface{}{1, 2}, []interface{}{"a", "b"}},
			expected: "",
		},
		{
			testName: "Resource",
			input:    file,
			expected: "",
		},
		{
			testName: "Custom type",
			input:    CustomType{"Hello world"},
			expected: "",
		},
	}
	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := zendParseArgAsString(tc.input)
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
