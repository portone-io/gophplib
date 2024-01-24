package gophplib

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleBase64Encode() {
	// Plain string
	fmt.Println(Base64Encode("Hello"))
	// Int
	fmt.Println(Base64Encode(123))
	// Float
	fmt.Println(Base64Encode(10.5))
	// Nil
	fmt.Println(Base64Encode(nil))
	// Empty string
	fmt.Println(Base64Encode(""))
	// Empty array
	fmt.Println(Base64Encode([]int{}))
	// Complex string with upper-case letters, numbers, and symbols
	fmt.Println(Base64Encode("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!%^&*(){}[]"))
	// String with control characters
	fmt.Println(Base64Encode("\n\t Line with control characters\r\n"))
	// String with non-standard characters
	fmt.Println(Base64Encode("\xC1\xC2\xC3\xC4\xC5\xC6"))
	// 한글
	fmt.Println(Base64Encode("안녕하세요"))

	// Output:
	// SGVsbG8= <nil>
	// MTIz <nil>
	// MTAuNQ== <nil>
	//  <nil>
	//  <nil>
	//  unsupported type : []int
	// QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVoxMjM0NTY3ODkwISVeJiooKXt9W10= <nil>
	// CgkgTGluZSB3aXRoIGNvbnRyb2wgY2hhcmFjdGVycw0K <nil>
	// wcLDxMXG <nil>
	// 7JWI64WV7ZWY7IS47JqU <nil>
}

func TestBase64Encode(t *testing.T) {
	testCases := []struct {
		any
		string
	}{
		{
			"Hello",
			"SGVsbG8=",
		},
		{
			123,
			"MTIz",
		},
		{
			10.5,
			"MTAuNQ==",
		},
		{
			nil,
			"",
		},
		{
			"",
			"",
		},
		{
			[]int{},
			"",
		},
		{
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!%^&*(){}[]",
			"QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVoxMjM0NTY3ODkwISVeJiooKXt9W10=",
		},
		{
			"\n\t Line with control characters\r\n",
			"CgkgTGluZSB3aXRoIGNvbnRyb2wgY2hhcmFjdGVycw0K",
		},
		{
			"\xC1\xC2\xC3\xC4\xC5\xC6",
			"wcLDxMXG",
		},
		{
			"안녕하세요",
			"7JWI64WV7ZWY7IS47JqU",
		},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := Base64Encode(tc.any)
			if err != nil {
				stringErr := fmt.Errorf("unsupported type : %s", reflect.TypeOf(tc.any))
				if err.Error() != stringErr.Error() {
					t.Errorf("%s: string error : %s, got %s", testName, stringErr, err)
				}
			} else {
				if !reflect.DeepEqual(result, tc.string) {
					t.Errorf("%s: byte %v, got %v", testName, tc.string, result)
				}
			}
		})
	}
}
