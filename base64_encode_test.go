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
	// Empty array: Unsupported type
	fmt.Println(Base64Encode([]int{}))

	// Output:
	// SGVsbG8= <nil>
	// MTIz <nil>
	// MTAuNQ== <nil>
	//  <nil>
	//  <nil>
	//  unsupported type : []int
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
			"????>>>>",
			"Pz8/Pz4+Pj4=",
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

func TestHugeBase64Error(t *testing.T) {
	result, err := Base64Encode(string(make([]byte, 1610612734)))
	if err == nil || err.Error() != "string too long, maximum is 1610612733" {
		t.Errorf("Expected error, got (%v, %v)", result, err)
	}
}
