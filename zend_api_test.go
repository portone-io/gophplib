package gophplib

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestZendParseArgAsString(t *testing.T) {

	testCase := []struct {
		any
		string
	}{
		{"Hello world", "Hello world"},
		{"", ""},
		{"Line1\nLine2\tTab", "Line1\nLine2\tTab"},
		{123, "123"},
		{9223372036854775807, "9223372036854775807"},
		{-123, "-123"},
		{123.456, "123.456"},
		{123.456789012345678, "123.45678901235"},
		{10.1234567e10, "101234567000"},
		{math.NaN(), "NAN"},
		{math.Inf(1), "INF"},
		{math.Inf(-1), "-INF"},
		{10.1234567e10, "101234567000"},
		{-123.456, "-123.456"},
		{0.0, "0"},
		{true, "1"},
		{false, ""},
		{nil, ""},
		{Sample{}, "sample object"},
		{Sample2{}, ""},
		{[]int{1, 2, 3}, ""},
		{[]string{"hello", "world"}, ""},
		{[]interface{}{[]interface{}{1, 2}, []interface{}{"a", "b"}}, ""},
		{getFile(), ""},
		{CustomType{"Hello world"}, ""},
	}
	for _, tc := range testCase {
		testName := fmt.Sprintf("%v", tc.any)
		t.Run(testName, func(t *testing.T) {
			result, err := zendParseArgAsString(tc.any)
			if err != nil {
				expectedErr := fmt.Errorf("unsupported type : %s", reflect.TypeOf(tc.any))
				if err.Error() != expectedErr.Error() {
					t.Errorf("%s: expected error : %s, got %s", testName, expectedErr, err)
				}
			} else {
				if !reflect.DeepEqual(result, tc.string) {
					t.Errorf("%s: expected %v, got %v", testName, tc.string, result)
				}
			}
		})
	}
}
