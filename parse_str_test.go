package gophplib

import (
	"fmt"
	"reflect"
	"testing"
)

// Basic test cases
func ExampleParseStr() {
	// Plain key-value
	fmt.Println(ParseStr("key=value&foo=bar"))

	// Array will be parsed as map with integer keys
	fmt.Println(ParseStr("arr[0]=A&arr[1]=B&arr[2]=C"))

	// Empty key will be treated as auto-incremented integer key for each array
	fmt.Println(ParseStr("arr[]=A&arr[]=B&arr[]=C&another[]=A&another[]=B"))

	// Dictionary
	fmt.Println(ParseStr("dict[key]=value&dict[foo]=bar"))

	// Nesting is allowed
	fmt.Println(ParseStr("dict[k1][k2]=v1&dict[k1][k3]=v2"))

	// ParseStr will automatically urldecode the input
	fmt.Println(ParseStr("firstname=Conan&surname=O%27Brien"))
	// Output:
	// map[foo:bar key:value]
	// map[arr:map[0:A 1:B 2:C]]
	// map[another:map[0:A 1:B] arr:map[0:A 1:B 2:C]]
	// map[dict:map[foo:bar key:value]]
	// map[dict:map[k1:map[k2:v1 k3:v2]]]
	// map[firstname:Conan surname:O'Brien]
}

// Test cases for complex array
func ExampleParseStr_complexArray() {
	// Each empty key will be treated as auto-incremented integer key for each
	// array
	fmt.Println(ParseStr("key=value&a[]=123&a[]=false&b[]=str&c[]=3.5&a[]=last"))

	// You can mix multiple types of keys in one dictionary, and you can mix
	// empty key with non-empty key. Each non-empty integer key will be used as
	// a new starting number for next empty key.
	fmt.Println(ParseStr("arr[]=A&arr[]=B&arr[9]=C&arr[]=D&arr[foo]=E&arr[]=F&arr[15.1]=G&arr[]=H"))

	// You can use empty key for multi-dimensional array. Refer to the following
	// example for the exact behavior.
	fmt.Println()
	fmt.Println("2-dim array:         ", ParseStr("arr[3][4]=deedee&arr[3][6]=wiz"))
	fmt.Println("2-dim with empty key:", ParseStr("arr[][]=deedee&arr[][]=wiz"))
	fmt.Println("partial empty key 1: ", ParseStr("arr[2][]=deedee&arr[2][]=wiz"))
	fmt.Println("partial empty key 2: ", ParseStr("arr[2][]=deedee&arr[4][]=wiz"))
	fmt.Println("partial empty key 3: ", ParseStr("arr[2][]=deedee&arr[][4]=wiz"))
	fmt.Println("partial empty key 4: ", ParseStr("arr[2][]=deedee&arr[][]=wiz"))
	fmt.Println("2-dim dict:          ", ParseStr("arr[one][four]=deedee&arr[three][six]=wiz"))
	fmt.Println("3-dim arr:           ", ParseStr("arr[1][2][3]=deedee&arr[1][2][6]=wiz"))
	// Output:
	// map[a:map[0:123 1:false 2:last] b:map[0:str] c:map[0:3.5] key:value]
	// map[arr:map[0:A 1:B 9:C 10:D 11:F 12:H 15.1:G foo:E]]
	//
	// 2-dim array:          map[arr:map[3:map[4:deedee 6:wiz]]]
	// 2-dim with empty key: map[arr:map[0:map[0:deedee] 1:map[0:wiz]]]
	// partial empty key 1:  map[arr:map[2:map[0:deedee 1:wiz]]]
	// partial empty key 2:  map[arr:map[2:map[0:deedee] 4:map[0:wiz]]]
	// partial empty key 3:  map[arr:map[2:map[0:deedee] 3:map[4:wiz]]]
	// partial empty key 4:  map[arr:map[2:map[0:deedee] 3:map[0:wiz]]]
	// 2-dim dict:           map[arr:map[one:map[four:deedee] three:map[six:wiz]]]
	// 3-dim arr:            map[arr:map[1:map[2:map[3:deedee 6:wiz]]]]
}

// Notable corner cases
func ExampleParseStr_cornerCases() {
	// input without key name will be ignored
	fmt.Println("empty input:", ParseStr(""))
	fmt.Println("no name:    ", ParseStr("=123&[]=123&[foo]=123&[3][var]=123"))
	fmt.Println("no value:   ", ParseStr("foo&arr[]&arr[]&arr[]=val"))
	fmt.Println()

	// ParseStr will automatically urldecode the input
	fmt.Println("encoded data:", ParseStr("a=%3c%3d%3d%20%20yolo+swag++%3d%3d%3e&b=%23%23%23Yolo+Swag%23%23%23"))
	fmt.Println("backslash:   ", ParseStr("sum=8%5c2%3d4"))
	fmt.Println("quotes:      ", ParseStr("str=%22quoted%22+string"))
	// Ill-formed urlencoded data will be ignored and remain unescaped
	fmt.Println("ill encoding:", ParseStr("first=%41&second=%a&third=%ZZ"))
	// Null bytes will be parsed as "%0"
	fmt.Println()

	// Some characters will be replaced with underscore
	fmt.Println("non-binary safe name:", ParseStr("arr.test[1]=deedee&arr test[4][two]=wiz"))
	fmt.Println("complex string:      ", ParseStr("first=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing"))
	fmt.Println("ill formed input:    ", ParseStr("yo;lo&foo = bar%ZZ&yolo + = + swag"))
	fmt.Println("ill formed key:      ", ParseStr("arr[1=deedee&arr[4][2=wiz"))
	// Output:
	// empty input: map[]
	// no name:     map[]
	// no value:    map[arr:map[0: 1: 2:val] foo:]
	//
	// encoded data: map[a:<==  yolo swag  ==> b:###Yolo Swag###]
	// backslash:    map[sum:8\2=4]
	// quotes:       map[str:"quoted" string]
	// ill encoding: map[first:A second:%a third:%ZZ]
	//
	// non-binary safe name: map[arr_test:map[1:deedee 4:map[two:wiz]]]
	// complex string:       map[arr:map[0:foo bar 1:baz] first:value foo:map[bar:foobar] test_field:testing]
	// ill formed input:     map[foo_: bar%ZZ yo;lo: yolo___:   swag]
	// ill formed key:       map[arr:map[4:wiz] arr_1:deedee]
}

func ExampleParseStr_version() {
	// parse_str("foo[ 3=v") returns ["foo_ 3" => "v"] in PHP 5.6 and
	// ["foo__3" => "v"] in PHP 8.3. We follows 5.6 behavior for compatibility.
	fmt.Println(ParseStr("foo[ 3=v"))
	// Output: map[foo_ 3:v]
}

// Test cases for ParseStr. These tests were created using the following test
// cases in PHP as inspiration.
//
// Reference:
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic1.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic2.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic3.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic4.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_memory_error.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/bug77439.phpt
//   - https://github.com/simnalamburt/snippets/blob/59843441/php/parse_str.php
func TestParseStr(t *testing.T) {
	type dict = map[any]any

	testCases := []struct {
		name     string
		input    string
		expected dict
	}{
		{
			name:  "BasicTest",
			input: "A=aaa&B=bbb&C=ccc",
			expected: dict{
				"A": "aaa",
				"B": "bbb",
				"C": "ccc",
			},
		},
		{
			name:  "ArrayValues",
			input: "A=aaa&a[]=111&a[]=true&b[]=bbb&c[]=1.414&a[]=3.14",
			expected: dict{
				"A": "aaa",
				"a": dict{0: "111", 1: "true", 2: "3.14"},
				"b": dict{0: "bbb"},
				"c": dict{0: "1.414"},
			},
		},
		{
			name:  "EncodedData",
			input: "a=%3c%3d%3d%20%20url+encoded++%3d%3d%3e&b=%23%23%23Url+Encoded%23%23%23",
			expected: dict{
				"a": "<==  url encoded  ==>",
				"b": "###Url Encoded###",
			},
		},
		{
			name:  "SingleQuotes",
			input: "firstname=Conan&surname=O%27Brien",
			expected: dict{
				"firstname": "Conan",
				"surname":   "O'Brien",
			},
		},
		{
			name:  "BackSlash",
			input: "sum=8%5c2%3d4",
			expected: dict{
				"sum": `8\2=4`,
			},
		},
		{
			name:  "DoubleQuotes",
			input: "str=%22quoted%22+string",
			expected: dict{
				"str": `"quoted" string`,
			},
		},
		{
			name:  "StringWithNulls",
			input: "str=string%20with%20%00%00%00%20nulls",
			expected: dict{
				"str": "string with \x00\x00\x00 nulls",
			},
		},
		{
			name:  "StringWith2DimArrayNumericKey",
			input: "arr[2][4]=deedee&arr[2][6]=wiz",
			expected: dict{
				"arr": dict{
					2: dict{
						4: "deedee",
						6: "wiz",
					},
				},
			},
		},
		{
			name:  "StringWith2DimArrayNullKey",
			input: "arr[][]=deedee&arr[][]=wiz",
			expected: dict{
				"arr": dict{
					0: dict{
						0: "deedee",
					},
					1: dict{
						0: "wiz",
					},
				},
			},
		},
		{
			name:  "StringWith2DimArrayNonNumericKey",
			input: "arr[a][d]=deedee&arr[c][six]=wiz",
			expected: dict{
				"arr": dict{
					"a": dict{
						"d": "deedee",
					},
					"c": dict{
						"six": "wiz",
					},
				},
			},
		},
		{
			name:  "StringWith3DimArrayNumericKey",
			input: "arr[1][2][3]=deedee&arr[1][2][10]=wiz",
			expected: dict{
				"arr": dict{
					1: dict{
						2: dict{
							3:  "deedee",
							10: "wiz",
						},
					},
				},
			},
		}, {
			name:  "StringWithNumericalArrayKeys",
			input: "arr[1]=deedee&arr[4]=sonny",
			expected: dict{
				"arr": dict{
					1: "deedee",
					4: "sonny",
				},
			},
		},
		{
			name:  "StringWithAssociativeKeys",
			input: "arr[A]=deedee&arr[D]=sonny",
			expected: dict{
				"arr": dict{
					"A": "deedee",
					"D": "sonny",
				},
			},
		},
		{
			name:  "BadlyFormedStrings0",
			input: "arr[1=deedee",
			expected: dict{
				"arr_1": "deedee",
			},
		},
		{
			name:  "BadlyFormedStrings1",
			input: "arr[1=deedee&arr[3][2=wiz",
			expected: dict{
				"arr_1": "deedee",
				"arr": dict{
					3: "wiz",
				},
			},
		},
		{
			name:  "BadlyFormedStrings2",
			input: "arr1]=deedee&arr[3]2]=wiz",
			expected: dict{
				"arr1]": "deedee",
				"arr": dict{
					3: "wiz",
				},
			},
		},
		{
			name:  "BadlyFormedStrings3",
			input: "arr[a=deedee&arr[4][b=wiz",
			expected: dict{
				"arr_a": "deedee",
				"arr": dict{
					4: "wiz",
				},
			},
		},
		{
			name:  "EncodedNumbers",
			input: "A=%41&B=%a&C=%b",
			expected: dict{
				"A": "A",
				"B": "%a",
				"C": "%b",
			},
		},
		{
			name:  "NonBinarySafeName",
			input: "arr.test[1]=deedee&arr test[4][b]=wiz",
			expected: dict{
				"arr_test": dict{
					1: "deedee",
					4: dict{
						"b": "wiz",
					},
				},
			},
		},
		{
			name:  "ComplexString",
			input: "A=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing",
			expected: dict{
				"A":          "value",
				"arr":        dict{0: "foo bar", 1: "baz"},
				"foo":        dict{"bar": "foobar"},
				"test_field": "testing",
			},
		},
		{
			name:  "IllInputString",
			input: "yo;lo&foo = bar%ZZ&yolo + = + swag",
			expected: dict{
				"yo;lo":   "",
				"foo_":    " bar%ZZ",
				"yolo___": "   swag",
			},
		},
		{
			name:  "MixKeyType",
			input: "2=222&3.14=3.14&arr[123]=asdf&arr[3.14]=asdf",
			expected: dict{
				2:      "222",
				"3_14": "3.14",
				"arr": dict{
					123:    "asdf",
					"3.14": "asdf",
				},
			},
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: dict{},
		},
		{
			name:     "NoName",
			input:    "=123&[]=123&[foo]=123&[3][var]=123",
			expected: dict{},
		},
		{
			name:     "NoValue",
			input:    "foo&arr[]&arr[]&arr[]=val",
			expected: dict{"foo": "", "arr": dict{0: "", 1: "", 2: "val"}},
		},
		{
			name:     "ReadingZero",
			input:    "foo[03]=yo",
			expected: dict{"foo": dict{"03": "yo"}},
		},
		{
			name:     "Negative",
			input:    "foo[-3]=yo",
			expected: dict{"foo": dict{-3: "yo"}},
		},
		{
			name:     "LeadingWhitespaces",
			input:    "   foo=bar",
			expected: dict{"foo": "bar"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseStr(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf(`
expected  %#v
actual    %#v`, tc.expected, result)
			}
		})
	}
}

// TODO: fuzzing

// Microbenchmark for ParseStr. Command:
//
//	go test -run '^$' -bench '^BenchmarkParseStr$' -benchmem \
//	  -benchtime 3s -count 5 -memprofile=mem.pprof -cpuprofile=cpu.pprof
func BenchmarkParseStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseStr("2=222&3.14=3.14&arr[123]=asdf&arr[3.14]=asdf&yo;lo&foo = bar%ZZ&yolo + = + swag&A=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing&arr.test[1]=deedee&arr test[4][b]=wiz")
	}
}

func TestZendHandleNumericStr(t *testing.T) {
	cases := []struct {
		string
		any
	}{
		{"", ""},
		{"0", 0},
		{"00", "00"},
		{"-0", "-0"},
		{"+0", "+0"},
		{"3", 3},
		{"03", "03"},
		{"-3", -3},
		{"+3", "+3"},
		{"3.14", "3.14"},
		{"314", 314},
		{"9999999999999999999999999999999", "9999999999999999999999999999999"},
		{"arr", "arr"},
	}

	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			if result := zendHandleNumericStr([]byte(c.string)); result != c.any {
				t.Errorf("expected %v, got %v", c.any, result)
			}
		})
	}
}
