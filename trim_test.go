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
		any
		string
	}{
		{"Hello world ", "Hello world"},
		{"", ""},
		{123, "123"},
		{-123, "-123"},
		{0, "0"},
		{true, "1"},
		{false, ""},
		{nil, ""},
		{123.40, "123.4"},
		{-123.40, "-123.4"},
		{10.1234567e10, "101234567000"},
		{1230.12984732591475609346509132875091237, "1230.1298473259"},
		{1230.12984732500000000000000000000000000, "1230.129847325"},
		{123456789123456.40, "1.2345678912346E+14"},
		{12345678912340.40, "12345678912340"},
		{c, c.toString()},
		{customTrim(nil), "hello world"},
		{`<header>
	<h1>hello world   </h1>
</header>`, "<header>\n\t<h1>hello world   </h1>\n</header>"},
		{" \x00\t\nABC \x00\t\n", "ABC"},
	}

	for _, tc := range testCase {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := Trim(tc.any)
			if err != nil {
				t.Errorf("%s: string success, bug got error %v", testName, err)
			}
			if !reflect.DeepEqual(result, tc.string) {
				t.Errorf("%s: string %v, bug got %v", testName, tc.string, result)
			}
		})
	}

	// Failing cases
	errorCase := []any{
		[]any{},
		[]any{"foo", "456", 7},
		d,
		file,
	}

	for _, tc := range errorCase {
		testName := fmt.Sprintf("%v", tc)
		t.Run(testName, func(t *testing.T) {
			result, err := Trim(tc)
			if err == nil {
				t.Errorf("%s:  error, bug got %v", testName, result)
			}
		})
	}
}
