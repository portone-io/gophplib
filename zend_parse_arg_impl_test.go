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

func ExampleZendParseArgImpl() {
	// Plain string
	fmt.Println(ZendParseArgImpl("Hello world"))

	// Empty string
	fmt.Println(ZendParseArgImpl(""))

	// String with special characters
	fmt.Println(ZendParseArgImpl("Line1\\nLine2\\tTab"))

	// Int
	var myInt int = 123
	fmt.Println(ZendParseArgImpl(myInt))

	// Int64
	var myInt64 int64 = 9223372036854775807
	fmt.Println(ZendParseArgImpl(myInt64))

	// Negative int
	myInt = -123
	fmt.Println(ZendParseArgImpl(myInt))

	// Float64
	var myFloat64 float64 = 123.456
	fmt.Println(ZendParseArgImpl(myFloat64))

	// Float64 exceeds 14 digits
	myFloat64 = 123.456789012345678
	fmt.Println(ZendParseArgImpl(myFloat64))

	// Exponent
	var myExponent = 10.1234567e10
	fmt.Println(ZendParseArgImpl(myExponent))

	// Special float - NaN
	var myNan = math.NaN()
	fmt.Println(ZendParseArgImpl(myNan))

	// Special float - positive infinity
	var myInf = math.Inf(1)
	fmt.Println(ZendParseArgImpl(myInf))

	// Special float - negative infinity
	myInf = math.Inf(-1)
	fmt.Println(ZendParseArgImpl(myInf))

	// Negative float64
	myFloat64 = -123.456
	fmt.Println(ZendParseArgImpl(myFloat64))

	// Zero float
	myFloat64 = 0.0
	fmt.Println(ZendParseArgImpl(myFloat64))

	// True
	var myTrue = true
	fmt.Println(ZendParseArgImpl(myTrue))

	// False
	var myFalse = false
	fmt.Println(ZendParseArgImpl(myFalse))

	// Nil
	fmt.Println(ZendParseArgImpl(nil))

	// Object with toString function
	var myObject = Sample{}
	fmt.Println(ZendParseArgImpl(myObject))

	// Object without toString function
	var myObject2 = Sample2{}
	fmt.Println(ZendParseArgImpl(myObject2))

	// Int Array
	var myArray []int = []int{1, 2, 3}
	fmt.Println(ZendParseArgImpl(myArray))

	// String array
	var myArray2 = []string{"hello", "world"}
	fmt.Println(ZendParseArgImpl(myArray2))

	// Resource
	file, osErr := os.Open("README.md")
	if osErr != nil {
		panic(osErr)
	}
	fmt.Println(ZendParseArgImpl(file))

	// Custom type
	var myCustom CustomType = CustomType{"Hello world"}
	fmt.Println(ZendParseArgImpl(myCustom))

	// Output:
	// Hello world <nil>
	//  <nil>
	// Line1\nLine2\tTab <nil>
	// 123 <nil>
	// 9223372036854775807 <nil>
	// -123 <nil>
	// 123.456 <nil>
	// 123.45678901235 <nil>
	// 101234567000 <nil>
	// NAN <nil>
	// INF <nil>
	// -INF <nil>
	// -123.456 <nil>
	// 0 <nil>
	// 1 <nil>
	//  <nil>
	//  <nil>
	// sample object <nil>
	// <nil> unsupported type : gophplib.Sample2
	// <nil> unsupported type : []int
	// <nil> unsupported type : []string
	// <nil> unsupported type : *os.File
	// <nil> unsupported type : gophplib.CustomType
}

func TestZendParseArgImpl(t *testing.T) {
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
			result, err := ZendParseArgImpl(tc.input)
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
