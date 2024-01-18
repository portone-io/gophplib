package gophplib

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type Cat struct {
	name string
	age  int
}

type Dog struct {
	name string
	age  int
}

func (c Cat) toString() string {
	return fmt.Sprintf("name is %s and %d years old", c.name, c.age)
}

func customTrim(value any) string {
	return " hello world "
}

func ExampleTrim_string() {
	// Trim plain string
	fmt.Println(Trim("Hello world "))

	// Trim empty string
	fmt.Println(Trim(""))

	// Trim string with default characters
	fmt.Println(Trim(" \x00\t\nABC \x00\t\n"))

	// Output:
	// Hello world <nil>
	//  <nil>
	// ABC <nil>
}

func ExampleTrim_float() {
	// Trim float
	fmt.Println(Trim(123.40))

	// Trim negative float
	fmt.Println(Trim(-123.40))

	// Trim exponent
	fmt.Println(Trim(10.1234567e10))

	// Trim float exceeds 14 digits
	fmt.Println(Trim(1230.12984732591475609346509132875091237))

	// Trim float ends with 0
	fmt.Println(Trim(1230.12984732500000000000000000000000000))

	// Trim integer exceeds 14 digits
	fmt.Println(Trim(123456789123456.40))

	// Trim integer ends with 0
	fmt.Println(Trim(12345678912340.40))

	// Output:
	// 123.4 <nil>
	// -123.4 <nil>
	// 101234567000 <nil>
	// 1230.1298473259 <nil>
	// 1230.129847325 <nil>
	// 1.2345678912346E+14 <nil>
	// 12345678912340 <nil>
}

func ExampleTrim_otherTypes() {
	// Trim integer
	fmt.Println(Trim(123))

	// Trim negative integer
	fmt.Println(Trim("-123"))

	// Trim zero
	fmt.Println(Trim(0))

	// Trim empty array
	fmt.Println(Trim([]any{}))

	// Trim bool (true)
	fmt.Println(Trim(true))

	// Trim bool (false)
	fmt.Println(Trim(false))

	// Trim object has toString
	fmt.Println(Trim(Cat{name: "nabi", age: 3}))

	// Trim object has no toString
	fmt.Println(Trim(Dog{name: "choco", age: 5}))

	// Trim function
	fmt.Println(Trim(customTrim(nil)))

	// Trim here documents
	fmt.Println(Trim(`<header>
	<h1>hello world   </h1>
</header>`))

	// Trim resource
	file, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	fmt.Println(Trim(file))

	// Output:
	// 123 <nil>
	// -123 <nil>
	// 0 <nil>
	//  unsupported type : []interface {}
	// 1 <nil>
	//  <nil>
	// name is nabi and 3 years old <nil>
	//  unsupported type : gophplib.Dog
	// hello world <nil>
	// <header>
	//	<h1>hello world   </h1>
	//</header> <nil>
	//  unsupported type : *os.File
}

func TestTrim(t *testing.T) {
	c := Cat{
		name: "nabi",
		age:  3,
	}

	d := Dog{
		name: "Choco",
		age:  4,
	}
	file, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}

	// Successful cases
	testCase := []struct {
		testName string
		input    any
		expected any
	}{
		{
			testName: "BasicTest",
			input:    "Hello world ",
			expected: "Hello world",
		},
		{
			testName: "EmptyString",
			input:    "",
			expected: "",
		},
		{
			testName: "Integer",
			input:    123,
			expected: "123",
		},
		{
			testName: "NegativeInteger",
			input:    -123,
			expected: "-123",
		},
		{
			testName: "Zero",
			input:    0,
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
			testName: "Float",
			input:    123.40,
			expected: "123.4",
		},
		{
			testName: "NegativeFloat",
			input:    -123.40,
			expected: "-123.4",
		},
		{
			testName: "Exponent",
			input:    10.1234567e10,
			expected: "101234567000",
		},
		{
			testName: "FloatExceeds14Digits",
			input:    1230.12984732591475609346509132875091237,
			expected: "1230.1298473259",
		},
		{
			testName: "FloatEndsWith0",
			input:    1230.12984732500000000000000000000000000,
			expected: "1230.129847325",
		},
		{
			testName: "IntegerExceeds14Digits",
			input:    123456789123456.40,
			expected: "1.2345678912346E+14",
		},
		{
			testName: "IntegerEndsWith0",
			input:    12345678912340.40,
			expected: "12345678912340",
		},
		{
			testName: "ObjectHasToString",
			input:    c,
			expected: c.toString(),
		},
		{
			testName: "Function",
			input:    customTrim(nil),
			expected: "hello world",
		},
		{
			testName: "HereDocuments",
			input: `<header>
	<h1>hello world   </h1>
</header>`,
			expected: "<header>\n\t<h1>hello world   </h1>\n</header>",
		},
		{
			testName: "StringWithDefaultCharacters",
			input:    " \x00\t\nABC \x00\t\n",
			expected: "ABC",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := Trim(tc.input)
			if err != nil {
				t.Errorf("%s: expected success, bug got error %v", tc.testName, err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("%s: expected %v, bug got %v", tc.testName, tc.expected, result)
			}
		})
	}

	// Failing cases
	errorCase := []struct {
		testName string
		input    any
		expected any
	}{
		{
			testName: "EmptyArray",
			input:    []any{},
		},
		{
			testName: "Array",
			input:    []any{"foo", "456", 7},
		},
		{
			testName: "ObjectWithoutToString",
			input:    d,
		},
		{
			testName: "Resource",
			input:    file,
		},
	}

	for _, tc := range errorCase {
		t.Run(tc.testName, func(t *testing.T) {
			result, err := Trim(tc.input)
			if err == nil {
				t.Errorf("%s: expected error, bug got %v", tc.testName, result)
			}
		})
	}
}
