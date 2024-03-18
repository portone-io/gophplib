package gophplib

import (
	"fmt"
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

func getFile() *os.File {
	file, osErr := os.Open("README.md")
	if osErr != nil {
		panic(osErr)
	}
	return file
}

func ExampleConvertToString() {
	// Nil
	fmt.Println(ConvertToString(nil))
	// Plain string
	fmt.Println(ConvertToString("Hello, World"))
	// Empty string
	fmt.Println(ConvertToString(""))
	// Int
	fmt.Println(ConvertToString(123))
	// Float
	fmt.Println(ConvertToString(123.45))
	// Exponent
	fmt.Println(ConvertToString(10.1234567e10))
	// Float exceeds 14 digits
	fmt.Println(ConvertToString(1230.12984732591475609346509132875091237))
	// Float ends with 0
	fmt.Println(ConvertToString(1230.12984732500000000000000000000000000))
	// True
	fmt.Println(ConvertToString(true))
	// False
	fmt.Println(ConvertToString(false))
	// Array
	fmt.Println(ConvertToString([]int{1, 2, 3}))
	// Slice
	fmt.Println(ConvertToString([2]int{1, 2}))
	// Object has toString
	fmt.Println(ConvertToString(Cat{
		name: "nabi",
		age:  3,
	}))
	// Object has no toString
	fmt.Println(ConvertToString(Dog{
		name: "choco",
		age:  5,
	}))

	// Output:
	//  <nil>
	// Hello, World <nil>
	//  <nil>
	// 123 <nil>
	// 123.45 <nil>
	// 101234567000 <nil>
	// 1230.1298473259 <nil>
	// 1230.129847325 <nil>
	// 1 <nil>
	//  <nil>
	// Array <nil>
	// Array <nil>
	// name is nabi and 3 years old <nil>
	//  unsupported type : gophplib.Dog
}

func TestConvertToString(t *testing.T) {
	file := getFile()
	testCase := []struct {
		any
		string
	}{
		{
			nil,
			"",
		},
		{
			"Hello, World",
			"Hello, World",
		},
		{
			"",
			"",
		},
		{
			123,
			"123",
		},
		{
			123.45,
			"123.45",
		},
		{
			10.1234567e10,
			"101234567000",
		},
		{
			1230.12984732591475609346509132875091237,
			"1230.1298473259",
		},
		{
			1230.12984732500000000000000000000000000,
			"1230.129847325",
		},
		{
			true,
			"1",
		},
		{
			false,
			"",
		},
		{
			[]int{1, 2, 3},
			"Array",
		},
		{
			[2]int{1, 2},
			"Array",
		},
		{
			file,
			fmt.Sprintf("Resource id %p", file),
		},
		{
			Cat{name: "nabi", age: 3},
			"name is nabi and 3 years old",
		},
	}

	for _, tc := range testCase {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := ConvertToString(tc.any)
			if err != nil {
				t.Errorf("%s: expected success to convert, bug get error %v", testName, err)

			}
			if !reflect.DeepEqual(result, tc.string) {
				t.Errorf("%s: string %v, bug got %v", testName, tc.string, result)
			}
		})
	}

	// Failing cases
	errorCase := []any{
		Dog{name: "choco", age: 5},
	}

	for _, tc := range errorCase {
		testName := fmt.Sprintf("%v", tc)
		t.Run(testName, func(t *testing.T) {
			result, err := ConvertToString(tc)
			if err == nil {
				t.Errorf("%s:  error, bug got %v", testName, result)
			}
		})
	}
	file.Close()
}
