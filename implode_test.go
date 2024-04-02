package gophplib

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
)

var multiTypedArr = []any{
	2,
	0,
	-639,
	true,
	"GO",
	false,
	nil,
	"",
	" ",
	"string\x00with\x00...\000",
}

func ExampleImplode() {
	// Arg1 is empty array and arg2 is nil
	fmt.Println(Implode([]any{}))

	// Arg1 is array and arg2 is nil
	fmt.Println(Implode([]string{"foo", "bar", "baz"}))

	// Arg1 is string and Arg2 is string array
	fmt.Println(Implode(":", []string{"foo", "bar", "baz"}))

	// Arg1 is string and arg2 is int array
	fmt.Println(Implode(", ", []int{1, 2}))

	// Arg1 is string and arg2 is float64 array
	fmt.Println(Implode(", ", []float64{1.1, 2.2}))

	// Arg1 is string and arg2 is boolean array
	fmt.Println(Implode(", ", []bool{false, true}))

	//Arg1 is string and arg2 is emtpy array
	fmt.Println(Implode(", ", []any{}))

	// Arg1 is string and Arg2 is 2D array
	fmt.Println(Implode(":", []any{"foo", []string{"bar", "baz"}, "burp"}))

	// Output:
	// <nil>
	// foobarbaz <nil>
	// foo:bar:baz <nil>
	// 1, 2 <nil>
	// 1.1, 2.2 <nil>
	// , 1 <nil>
	//  <nil>
	// foo:Array:burp <nil>
}

func ExampleImplode_variation() {
	// Arg1 is string 'TRUE' and arg2 is multi typed array
	result, ok := Implode("TRUE", multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is true and arg2 is multi typed array
	result, ok = Implode(true, multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is false and arg2 is multi typed array
	result, ok = Implode(false, multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is array and arg2 is multi typed array
	fmt.Println(Implode([]string{"key1", "key2"}, multiTypedArr))

	// Arg1 is empty string and arg2 is multi typed array
	result, ok = Implode("", multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is blank string and arg2 is multi typed array
	result, ok = Implode(" ", multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is string contains null bytes string and arg2 is multi typed array
	result, ok = Implode("bet\x00ween", multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is nil and arg2 is multi typed array
	result, ok = Implode(nil, multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is negative int arg2 is multi typed array
	result, ok = Implode(-0, multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is null bytes string arg2 is multi typed array
	result, ok = Implode(`\0`, multiTypedArr)
	fmt.Println(strings.Replace(result, "\x00", "NUL", -1), ok)

	// Arg1 is array and arg2 is not nil
	fmt.Println(Implode([]any{1, "2", 3.45, true}, "sep"))

	// Arg1 is string and arg2 is object array
	fmt.Println(Implode(", ", []any{Cat{"nabi", 3}}))

	// Initialize small scale map
	smallMap := orderedmap.NewOrderedMap[any, any]()
	smallMap.Set("key1", "value1")
	smallMap.Set("key2", "value2")
	fmt.Println(Implode(", ", smallMap))

	// Output:
	// 2TRUE0TRUE-639TRUE1TRUEGOTRUETRUETRUETRUE TRUEstringNULwithNUL...NUL <nil>
	// 2101-639111GO1111 1stringNULwithNUL...NUL <nil>
	// 20-6391GO stringNULwithNUL...NUL <nil>
	// key1Arraykey2 <nil>
	// 20-6391GO stringNULwithNUL...NUL <nil>
	// 2 0 -639 1 GO      stringNULwithNUL...NUL <nil>
	// 2betNULween0betNULween-639betNULween1betNULweenGObetNULweenbetNULweenbetNULweenbetNULween betNULweenstringNULwithNUL...NUL <nil>
	// 20-6391GO stringNULwithNUL...NUL <nil>
	// 2000-639010GO0000 0stringNULwithNUL...NUL <nil>
	// 2\00\0-639\01\0GO\0\0\0\0 \0stringNULwithNUL...NUL <nil>
	// 1sep2sep3.45sep1 <nil>
	// name is nabi and 3 years old <nil>
	// value1, value2 <nil>
}

func ExampleImplode_error() {
	// File resource
	fmt.Println(Implode(", ", getFile()))
	// Only arg1
	fmt.Println(Implode("glue"))

	// Arg2 is int
	fmt.Println(Implode("glue", 1234))

	// Arg2 is nil
	fmt.Println(Implode("glue", nil))

	// Arg1 is int
	fmt.Println(Implode(12, "pieces"))

	// Arg1 is nil
	fmt.Println(Implode(nil, "abcd"))

	// Output:
	//  invalid arguments passed, got string, *os.File
	//  argument must be one of array, slice, or ordered map, but got string
	//  invalid arguments passed, got string, int
	//  invalid arguments passed, got string, <nil>
	//  invalid arguments passed, got int, string
	//  invalid arguments passed, got <nil>, string
}

func TestImplode(t *testing.T) {
	om := orderedmap.NewOrderedMap[any, any]()
	om.Set("key1", "value1")
	om.Set("key2", "value2")
	testCases := []struct {
		arg1 any
		arg2 any
		string
	}{
		{[]any{}, nil, ""},
		{[]string{"foo", "bar", "baz"}, nil, "foobarbaz"},
		{":", []string{"foo", "bar", "baz"}, "foo:bar:baz"},
		{", ", []int{1, 2}, "1, 2"},
		{", ", []float64{1.1, 2.2}, "1.1, 2.2"},
		{", ", []bool{false, true}, ", 1"},
		{", ", []any{}, ""},
		{":", []any{"foo", []string{"bar", "baz"}, "burp"}, "foo:Array:burp"},
		{"TRUE", multiTypedArr, "2TRUE0TRUE-639TRUE1TRUEGOTRUETRUETRUETRUE TRUEstring\x00with\x00...\x00"},
		{true, multiTypedArr, "2101-639111GO1111 1string\x00with\x00...\x00"},
		{false, multiTypedArr, "20-6391GO string\x00with\x00...\x00"},
		{[]string{"key1", "key2"}, multiTypedArr, "key1Arraykey2"},
		{"", multiTypedArr, "20-6391GO string\x00with\x00...\x00"},
		{" ", multiTypedArr, "2 0 -639 1 GO      string\x00with\x00...\x00"},
		{"bet\x00ween", multiTypedArr, "2bet\x00ween0bet\x00ween-639bet\x00ween1bet\x00weenGObet\x00weenbet\x00weenbet\u0000weenbet\u0000ween bet\x00weenstring\x00with\x00...\x00"},
		{nil, multiTypedArr, "20-6391GO string\x00with\x00...\x00"},
		{-0, multiTypedArr, "2000-639010GO0000 0string\x00with\x00...\x00"},
		{`\0`, multiTypedArr, "2\\00\\0-639\\01\\0GO\\0\\0\\0\\0 \\0string\x00with\x00...\x00"},
		{[]any{1, "2", 3.45, true}, "sep", "1sep2sep3.45sep1"},
		{"glue", nil, ""},
		{"glue", 1234, ""},
		{"glue", nil, ""},
		{12, "pieces", ""},
		{nil, "abcd", ""},
		{", ", []any{Cat{"nabi", 3}}, "name is nabi and 3 years old"},
		{", ", map[string]string{"foo": "bar"}, "bar"},
		{", ", om, "value1, value2"},
		{", ", *om, "value1, value2"},
		{", ", []any{Dog{"choco", 5}, Cat{"nabi", 3}}, "Object, name is nabi and 3 years old"},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("arg 1 : %v, arg2 : %v", tc.arg1, tc.arg2)
		t.Run(testName, func(t *testing.T) {
			if tc.arg2 == nil {
				result, err := Implode(tc.arg1)
				if err != nil {
					if tc.string != "" {
						t.Errorf("%s: expected : %s, got error %s", testName, tc.string, err.Error())
					} else {
						expectedErr := fmt.Errorf("argument must be one of array, slice, or ordered map, but got %v", reflect.TypeOf(tc.arg1))
						if err.Error() != expectedErr.Error() {
							t.Errorf("%s: expected error : %s, got %s", testName, expectedErr.Error(), err.Error())
						}
					}
				} else {
					if !reflect.DeepEqual(result, tc.string) {
						t.Errorf("%s: expected : %s, got %s", testName, tc.string, result)
					}
				}
			} else {
				result, err := Implode(tc.arg1, tc.arg2)
				if err != nil {
					if tc.string != "" {
						t.Errorf("%s: expected : %s, got error %s", testName, tc.string, err.Error())
					} else {
						expectedErr := fmt.Errorf("invalid arguments passed, got %v, %v", reflect.TypeOf(tc.arg1), reflect.TypeOf(tc.arg2))
						if err.Error() != expectedErr.Error() {
							t.Errorf("%s: expected error : %s, got %s", testName, expectedErr.Error(), err.Error())
						}
					}
				} else {
					if !reflect.DeepEqual(result, tc.string) {
						t.Errorf("%s: expected : %s, got %s", testName, tc.string, result)
					}
				}
			}
		})
	}

	typeErrCase := struct {
		arg1 any
		arg2 any
	}{"foo", []any{1 + 2i, Cat{"nabi", 3}}}

	testName := fmt.Sprintf("%v", typeErrCase)
	t.Run(testName, func(t *testing.T) {
		result, err := Implode(typeErrCase.arg1, typeErrCase.arg2)
		if err != nil {
			if !strings.Contains(err.Error(), "unsupported type in array") {
				t.Errorf("%s: expected error : unsupported type in array, bug got %s", testName, err.Error())
			}
		} else {
			t.Errorf("%s: error, but got %v", testName, result)
		}
	})
}

func BenchmarkImplode(b *testing.B) {
	// Initialization
	var (
		// Small scale map
		smallMap = orderedmap.NewOrderedMap[any, any]()

		// Medium scale map
		mediumMap = orderedmap.NewOrderedMap[any, any]()

		// Large scale map
		largeMap = orderedmap.NewOrderedMap[any, any]()

		// String-based map
		stringMap = orderedmap.NewOrderedMap[any, any]()

		// Integer-based map
		intMap = orderedmap.NewOrderedMap[any, any]()
	)

	// Initialize small scale map
	smallMap.Set("key1", "value1")
	smallMap.Set("key2", "value2")

	// Initialize medium scale map
	for i := 0; i < 50; i++ {
		mediumMap.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// Initialize large scale map
	for i := 0; i < 1000; i++ {
		largeMap.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// Initialize string-based map
	for i := 0; i < 100; i++ {
		stringMap.Set(fmt.Sprintf("str%d", i), fmt.Sprintf("val%d", i))
	}

	// Initialize integer-based map
	for i := 0; i < 100; i++ {
		intMap.Set(i, i)
	}

	testCases := []struct {
		name string
		arg1 any
		arg2 any
	}{
		{"SmallMap", smallMap, nil},
		{"MediumMap", mediumMap, nil},
		{"LargeMap", largeMap, nil},
		{"StringKeysMap", stringMap, nil},
		{"IntegerKeysMap", intMap, nil},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var arg2Any []any
			if tc.arg2 != nil {
				switch v := tc.arg2.(type) {
				case string:
					arg2Any = append(arg2Any, v)
				case []string:
					for _, item := range v {
						arg2Any = append(arg2Any, item)
					}
				default:
					arg2Any = append(arg2Any, v)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if tc.arg2 == nil {
					_, _ = Implode(tc.arg1)
				} else {
					_, _ = Implode(tc.arg1, arg2Any...)
				}
			}
		})
	}
}

func TestIsOrderedMap(t *testing.T) {
	testCases := []struct {
		any
	}{
		{orderedmap.NewOrderedMap[any, any]()},
		{*orderedmap.NewOrderedMap[any, any]()},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%v", reflect.TypeOf(tc.any).Kind())
		t.Run(testName, func(t *testing.T) {
			if !isOrderedMap(tc.any) {
				t.Errorf("expected result was true but got false")
			}
		})
	}
}
