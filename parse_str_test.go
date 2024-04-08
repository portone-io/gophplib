package gophplib

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
)

func dumpOrderedMap(omap orderedmap.OrderedMap[any, any]) string {
	if omap.Len() <= 0 {
		return "omap[]"
	}

	builder := strings.Builder{}
	builder.WriteString("omap")

	first := true

	for el := omap.Front(); el != nil; el = el.Next() {
		if first {
			builder.WriteRune('[')
			first = false
		} else {
			builder.WriteRune(' ')
		}

		builder.WriteString(fmt.Sprint(el.Key))
		builder.WriteString(":")

		switch v := el.Value.(type) {
		case orderedmap.OrderedMap[any, any]:
			builder.WriteString(dumpOrderedMap(v))
		default:
			builder.WriteString(fmt.Sprint(el.Value))
		}
	}

	builder.WriteRune(']')
	return builder.String()
}

func omap(args ...any) orderedmap.OrderedMap[any, any] {
	omap := *orderedmap.NewOrderedMap[any, any]()

	var lastkey any
	lastkey = nil

	for _, v := range args {
		if lastkey == nil {
			lastkey = v
		} else {
			omap.Set(lastkey, v)
			lastkey = nil
		}
	}

	return omap
}

// Basic test cases
func ExampleParseStr() {
	// Plain key-value
	fmt.Println(dumpOrderedMap(ParseStr("key=value&foo=bar")))

	// Array will be parsed as map with integer keys
	fmt.Println(dumpOrderedMap(ParseStr("arr[0]=A&arr[1]=B&arr[2]=C")))

	// Empty key will be treated as auto-incremented integer key for each array
	fmt.Println(dumpOrderedMap(ParseStr("arr[]=A&arr[]=B&arr[]=C&another[]=A&another[]=B")))

	// Dictionary
	fmt.Println(dumpOrderedMap(ParseStr("dict[key]=value&dict[foo]=bar")))

	// Nesting is allowed
	fmt.Println(dumpOrderedMap(ParseStr("dict[k1][k2]=v1&dict[k1][k3]=v2")))

	// ParseStr will automatically urldecode the input
	fmt.Println(dumpOrderedMap(ParseStr("firstname=Conan&surname=O%27Brien")))

	// Output:
	// omap[key:value foo:bar]
	// omap[arr:omap[0:A 1:B 2:C]]
	// omap[arr:omap[0:A 1:B 2:C] another:omap[0:A 1:B]]
	// omap[dict:omap[key:value foo:bar]]
	// omap[dict:omap[k1:omap[k2:v1 k3:v2]]]
	// omap[firstname:Conan surname:O'Brien]
}

// Test cases for complex arrays.
//
// NOTE:
// Noticed that when using fmt.Println with a map, the output may vary depending on the version of Go being used.
// In certain cases, the output differed between versions, so those cases were commented out.
// However, please note that this discrepancy occurs only when using fmt.Println and not in the actual operation of the ParseStr function,
// which behaves as expected.
func ExampleParseStr_complexArray() {
	// Each empty key will be treated as auto-incremented integer key for each
	// array
	fmt.Println(dumpOrderedMap(ParseStr("key=value&a[]=123&a[]=false&b[]=str&c[]=3.5&a[]=last")))

	// You can mix multiple types of keys in one dictionary, and you can mix
	// empty key with non-empty key. Each non-empty integer key will be used as
	// a new starting number for next empty key.
	fmt.Println(dumpOrderedMap(ParseStr("arr[]=A&arr[]=B&arr[9]=C&arr[]=D&arr[foo]=E&arr[]=F&arr[15.1]=G&arr[]=H")))

	// You can use empty key for multi-dimensional array. Refer to the following
	// example for the exact behavior.
	fmt.Println("2-dim array:         ", dumpOrderedMap(ParseStr("arr[3][4]=deedee&arr[3][6]=wiz")))
	fmt.Println("2-dim with empty key:", dumpOrderedMap(ParseStr("arr[][]=deedee&arr[][]=wiz")))
	fmt.Println("partial empty key 1: ", dumpOrderedMap(ParseStr("arr[2][]=deedee&arr[2][]=wiz")))
	fmt.Println("partial empty key 2: ", dumpOrderedMap(ParseStr("arr[2][]=deedee&arr[4][]=wiz")))
	fmt.Println("partial empty key 3: ", dumpOrderedMap(ParseStr("arr[2][]=deedee&arr[][4]=wiz")))
	fmt.Println("partial empty key 4: ", dumpOrderedMap(ParseStr("arr[2][]=deedee&arr[][]=wiz")))
	fmt.Println("2-dim dict:          ", dumpOrderedMap(ParseStr("arr[one][four]=deedee&arr[three][six]=wiz")))
	fmt.Println("3-dim arr:           ", dumpOrderedMap(ParseStr("arr[1][2][3]=deedee&arr[1][2][6]=wiz")))

	// Output:
	// omap[key:value a:omap[0:123 1:false 2:last] b:omap[0:str] c:omap[0:3.5]]
	// omap[arr:omap[0:A 1:B 9:C 10:D foo:E 11:F 15.1:G 12:H]]
	// 2-dim array:          omap[arr:omap[3:omap[4:deedee 6:wiz]]]
	// 2-dim with empty key: omap[arr:omap[0:omap[0:deedee] 1:omap[0:wiz]]]
	// partial empty key 1:  omap[arr:omap[2:omap[0:deedee 1:wiz]]]
	// partial empty key 2:  omap[arr:omap[2:omap[0:deedee] 4:omap[0:wiz]]]
	// partial empty key 3:  omap[arr:omap[2:omap[0:deedee] 3:omap[4:wiz]]]
	// partial empty key 4:  omap[arr:omap[2:omap[0:deedee] 3:omap[0:wiz]]]
	// 2-dim dict:           omap[arr:omap[one:omap[four:deedee] three:omap[six:wiz]]]
	// 3-dim arr:            omap[arr:omap[1:omap[2:omap[3:deedee 6:wiz]]]]
}

// Notable corner cases
func ExampleParseStr_cornerCases() {
	// input without key name will be ignored
	fmt.Println("empty input:", dumpOrderedMap(ParseStr("")))
	fmt.Println("no name:    ", dumpOrderedMap(ParseStr("=123&[]=123&[foo]=123&[3][var]=123")))
	fmt.Println("no value:   ", dumpOrderedMap(ParseStr("foo&arr[]&arr[]&arr[]=val")))

	// ParseStr will automatically urldecode the input
	fmt.Println("encoded data:", dumpOrderedMap(ParseStr("a=%3c%3d%3d%20%20yolo+swag++%3d%3d%3e&b=%23%23%23Yolo+Swag%23%23%23")))
	fmt.Println("backslash:   ", dumpOrderedMap(ParseStr("sum=8%5c2%3d4")))
	fmt.Println("quotes:      ", dumpOrderedMap(ParseStr("str=%22quoted%22+string")))

	// Ill-formed urlencoded data will be ignored and remain unescaped
	fmt.Println("ill encoding:", dumpOrderedMap(ParseStr("first=%41&second=%a&third=%ZZ")))
	// Null bytes will be parsed as "%0"

	// Some characters will be replaced with underscore
	fmt.Println("non-binary safe name:", dumpOrderedMap(ParseStr("arr.test[1]=deedee&arr test[4][two]=wiz")))
	fmt.Println("complex string:      ", dumpOrderedMap(ParseStr("first=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing")))
	fmt.Println("ill formed input:    ", dumpOrderedMap(ParseStr("yo;lo&foo = bar%ZZ&yolo + = + swag")))
	fmt.Println("ill formed key:      ", dumpOrderedMap(ParseStr("arr[1=deedee&arr[4][2=wiz")))

	// Output:
	// empty input: omap[]
	// no name:     omap[]
	// no value:    omap[foo: arr:omap[0: 1: 2:val]]
	// encoded data: omap[a:<==  yolo swag  ==> b:###Yolo Swag###]
	// backslash:    omap[sum:8\2=4]
	// quotes:       omap[str:"quoted" string]
	// ill encoding: omap[first:A second:%a third:%ZZ]
	// non-binary safe name: omap[arr_test:omap[1:deedee 4:omap[two:wiz]]]
	// complex string:       omap[first:value arr:omap[0:foo bar 1:baz] foo:omap[bar:foobar] test_field:testing]
	// ill formed input:     omap[yo;lo: foo_: bar%ZZ yolo___:   swag]
	// ill formed key:       omap[arr_1:deedee arr:omap[4:wiz]]
}

func ExampleParseStr_version() {
	// parse_str("foo[ 3=v") returns ["foo_ 3" => "v"] in PHP 5.6 and
	// ["foo__3" => "v"] in PHP 8.3. We follows 5.6 behavior for compatibility.
	fmt.Println(dumpOrderedMap(ParseStr("foo[ 3=v")))
	// Output: omap[foo_ 3:v]
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
	testCases := []struct {
		name     string
		input    string
		expected orderedmap.OrderedMap[any, any]
	}{
		{
			name:  "BasicTest",
			input: "A=aaa&B=bbb&C=ccc",
			expected: omap(
				"A", "aaa",
				"B", "bbb",
				"C", "ccc",
			),
		},
		{
			name:  "ArrayValues",
			input: "A=aaa&a[]=111&a[]=true&b[]=bbb&c[]=1.414&a[]=3.14",
			expected: omap(
				"A", "aaa",
				"a", omap(0, "111", 1, "true", 2, "3.14"),
				"b", omap(0, "bbb"),
				"c", omap(0, "1.414"),
			),
		},
		{
			name:  "EncodedData",
			input: "a=%3c%3d%3d%20%20url+encoded++%3d%3d%3e&b=%23%23%23Url+Encoded%23%23%23",
			expected: omap(
				"a", "<==  url encoded  ==>",
				"b", "###Url Encoded###",
			),
		},
		{
			name:  "SingleQuotes",
			input: "firstname=Conan&surname=O%27Brien",
			expected: omap(
				"firstname", "Conan",
				"surname", "O'Brien",
			),
		},
		{
			name:  "BackSlash",
			input: "sum=8%5c2%3d4",
			expected: omap(
				"sum", `8\2=4`,
			),
		},
		{
			name:  "DoubleQuotes",
			input: "str=%22quoted%22+string",
			expected: omap(
				"str", `"quoted" string`,
			),
		},
		{
			name:  "StringWithNulls",
			input: "str=string%20with%20%00%00%00%20nulls",
			expected: omap(
				"str", "string with \x00\x00\x00 nulls",
			),
		},
		{
			name:  "StringWith2DimArrayNumericKey",
			input: "arr[2][4]=deedee&arr[2][6]=wiz",
			expected: omap(
				"arr", omap(
					2, omap(
						4, "deedee",
						6, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayNullKey",
			input: "arr[][]=deedee&arr[][]=wiz",
			expected: omap(
				"arr", omap(
					0, omap(
						0, "deedee",
					),
					1, omap(
						0, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayNonNumericKey",
			input: "arr[a][d]=deedee&arr[c][six]=wiz",
			expected: omap(
				"arr", omap(
					"a", omap(
						"d", "deedee",
					),
					"c", omap(
						"six", "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayPartialEmptyKey",
			input: "arr[2][]=deedee&arr[2][]=wiz",
			expected: omap(
				"arr", omap(
					2, omap(
						0, "deedee",
						1, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayPartialEmptyKey2",
			input: "arr[2][]=deedee&arr[4][]=wiz",
			expected: omap(
				"arr", omap(
					2, omap(
						0, "deedee",
					),
					4, omap(
						0, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayPartialEmptyKey3",
			input: "arr[2][]=deedee&arr[][4]=wiz",
			expected: omap(
				"arr", omap(
					2, omap(
						0, "deedee",
					),
					3, omap(
						4, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith2DimArrayPartialEmptyKey4",
			input: "arr[2][]=deedee&arr[][]=wiz",
			expected: omap(
				"arr", omap(
					2, omap(
						0, "deedee",
					),
					3, omap(
						0, "wiz",
					),
				),
			),
		},
		{
			name:  "StringWith3DimArrayNumericKey",
			input: "arr[1][2][3]=deedee&arr[1][2][10]=wiz",
			expected: omap(
				"arr", omap(
					1, omap(
						2, omap(
							3, "deedee",
							10, "wiz",
						),
					),
				),
			),
		}, {
			name:  "StringWithNumericalArrayKeys",
			input: "arr[1]=deedee&arr[4]=sonny",
			expected: omap(
				"arr", omap(
					1, "deedee",
					4, "sonny",
				),
			),
		},
		{
			name:  "StringWithAssociativeKeys",
			input: "arr[A]=deedee&arr[D]=sonny",
			expected: omap(
				"arr", omap(
					"A", "deedee",
					"D", "sonny",
				),
			),
		},
		{
			name:  "BadlyFormedStrings0",
			input: "arr[1=deedee",
			expected: omap(
				"arr_1", "deedee",
			),
		},
		{
			name:  "BadlyFormedStrings1",
			input: "arr[1=deedee&arr[3][2=wiz",
			expected: omap(
				"arr_1", "deedee",
				"arr", omap(
					3, "wiz",
				),
			),
		},
		{
			name:  "BadlyFormedStrings2",
			input: "arr1]=deedee&arr[3]2]=wiz",
			expected: omap(
				"arr1]", "deedee",
				"arr", omap(
					3, "wiz",
				),
			),
		},
		{
			name:  "BadlyFormedStrings3",
			input: "arr[a=deedee&arr[4][b=wiz",
			expected: omap(
				"arr_a", "deedee",
				"arr", omap(
					4, "wiz",
				),
			),
		},
		{
			name:  "EncodedNumbers",
			input: "A=%41&B=%a&C=%b",
			expected: omap(
				"A", "A",
				"B", "%a",
				"C", "%b",
			),
		},
		{
			name:  "NonBinarySafeName",
			input: "arr.test[1]=deedee&arr test[4][b]=wiz",
			expected: omap(
				"arr_test", omap(
					1, "deedee",
					4, omap(
						"b", "wiz",
					),
				),
			),
		},
		{
			name:  "ComplexString",
			input: "A=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing",
			expected: omap(
				"A", "value",
				"arr", omap(0, "foo bar", 1, "baz"),
				"foo", omap("bar", "foobar"),
				"test_field", "testing",
			),
		},
		{
			name:  "MixMultipleType",
			input: "arr[]=A&arr[]=B&arr[9]=C&arr[]=D&arr[foo]=E&arr[]=F&arr[15.1]=G&arr[]=H",
			expected: omap(
				"arr", omap(
					0, "A",
					1, "B",
					9, "C",
					10, "D",
					"foo", "E",
					11, "F",
					"15.1", "G",
					12, "H",
				),
			),
		},
		{
			name:  "IllInputString",
			input: "yo;lo&foo = bar%ZZ&yolo + = + swag",
			expected: omap(
				"yo;lo", "",
				"foo_", " bar%ZZ",
				"yolo___", "   swag",
			),
		},
		{
			name:  "MixKeyType",
			input: "2=222&3.14=3.14&arr[123]=asdf&arr[3.14]=asdf",
			expected: omap(
				2, "222",
				"3_14", "3.14",
				"arr", omap(
					123, "asdf",
					"3.14", "asdf",
				),
			),
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: omap(),
		},
		{
			name:     "NoName",
			input:    "=123&[]=123&[foo]=123&[3][var]=123",
			expected: omap(),
		},
		{
			name:     "NoValue",
			input:    "foo&arr[]&arr[]&arr[]=val",
			expected: omap("foo", "", "arr", omap(0, "", 1, "", 2, "val")),
		},
		{
			name:     "ReadingZero",
			input:    "foo[03]=yo",
			expected: omap("foo", omap("03", "yo")),
		},
		{
			name:     "Negative",
			input:    "foo[-3]=yo",
			expected: omap("foo", omap(-3, "yo")),
		},
		{
			name:     "LeadingWhitespaces",
			input:    "   foo=bar",
			expected: omap("foo", "bar"),
		},
		{
			name:     "ArrayThenString",
			input:    "foo[]=x&foo[]=y&foo=bar",
			expected: omap("foo", "bar"),
		},
		{
			name:     "StringThenArray",
			input:    "foo=bar&foo[]=x&foo[]=y",
			expected: omap("foo", omap(0, "x", 1, "y")),
		},
		{
			name:     "StringHasBlanks",
			input:    "foo[ 3=v",
			expected: omap("foo_ 3", "v"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := dumpOrderedMap(ParseStr(tc.input))
			expected := dumpOrderedMap(tc.expected)
			if !reflect.DeepEqual(result, expected) {
				t.Errorf(`
expected  %#v
actual    %#v`, tc.expected, result)
			}
		})
	}
}

// Microbenchmark for ParseStr. Command:
//
//	go test -run '^$' -bench '^BenchmarkParseStr$' -benchmem \
//	  -benchtime 3s -count 5 -memprofile=mem.pprof -cpuprofile=cpu.pprof
func BenchmarkParseStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseStr("2=222&3.14=3.14&arr[123]=asdf&arr[3.14]=asdf&yo;lo&foo = bar%ZZ&yolo + = + swag&A=value&arr[]=foo+bar&arr[]=baz&foo[bar]=foobar&test.field=testing&arr.test[1]=deedee&arr test[4][b]=wiz")
	}
}

func TestPhpNumericOrString(t *testing.T) {
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
			if result := phpNumericOrString([]byte(c.string)); result != c.any {
				t.Errorf("expected %v, got %v", c.any, result)
			}
		})
	}
}
