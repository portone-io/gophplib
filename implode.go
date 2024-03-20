package gophplib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

// Implode replicates the behavior of PHP 5.6's implode function in GO.
// This function concatenates the elements of an array into a single string using a specified separator.
// For more information, see the [official PHP documentation].
//
// The function supports flexible argument patterns to accommodate various use cases.
//   - arg1: Can be either a string (used as the separator) or an array of elements to be joined.
//   - options: Optional. When provided, the first element is used as arg2.
//     If arg1 is a string, arg2 serves as the array of elements to be joined.
//     If arg1 is an array, arg2 serves as the separator.
//
// Behavior:
//
// 1. If only arg1 is provided:
//   - If arg1 is an array, the function joins the elements using an empty string as the default separator.
//   - If arg1 is not an array, the function returns an error.
//
// 2. If both arg1 and arg2 are provided:
//   - If arg1 is an array, arg2 is converted to string and used as the separator.
//   - If arg1 is not an array and arg2 is an array, arg1 is converted to a string
//     and used as the separator, with arg2 being the array to implode.
//   - If neither arg1 nor arg2 is an array, the function returns an error.
//
// Non-string elements within the array are converted to strings using a ConvertToString function
// before joining.
// Due to language differences between PHP and Go, the implode function support OrderedMap type from the [orderedmap library],
// ensuring ordered map functionality. When imploding map types, please utilize the OrderedMap type from the [orderedmap library]
// to maintain element order. If you use map type, not OrderedMap type, the order of the results cannot be guaranteed.
//
// reference:
//   - implode: https://github.com/php/php-src/blob/php-5.6.40/ext/standard/string.c#L1229-L1269
//   - php_implode: https://github.com/php/php-src/blob/php-5.6.40/ext/standard/string.c#L1141-L1224
//
// Test cases:
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/implode.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/implode1.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.implode.php
// [orderedmap library]: https://pkg.go.dev/github.com/elliotchance/orderedmap/v2
func Implode(arg1 any, options ...any) (string, error) {
	var delim string
	var arr []any

	// Check arg1 is one of array, slice, map, or ordered ap
	isArg1CollectionType := isCollectionType(arg1)

	// Check if options is not provided
	if len(options) == 0 {
		if !isArg1CollectionType {
			return "", fmt.Errorf("argument must be one of array, slice, or ordered map, but got %v", reflect.TypeOf(arg1))
		}
		arr = aggregateValues(arg1)
	} else {
		arg2 := options[0]
		// Check arg2 is one of array, slice, or ordered map
		isArg2CollectionType := isCollectionType(arg2)

		if isArg1CollectionType {
			delim, _ = ConvertToString(arg2)
			arr = aggregateValues(arg1)
		} else if !isArg1CollectionType && isArg2CollectionType {
			delim, _ = ConvertToString(arg1)
			arr = aggregateValues(arg2)
		} else {
			return "", fmt.Errorf("invalid arguments passed, got %v, %v", reflect.TypeOf(arg1), reflect.TypeOf(arg2))
		}
	}

	// Join arr elements with a delim
	var builder strings.Builder
	if len(arr) == 0 {
		return "", nil
	}
	for i, item := range arr {
		str, err := ConvertToString(item)

		if err != nil {
			return "", fmt.Errorf("unsupported type in array : %v", reflect.TypeOf(item))
		} else {
			builder.WriteString(str)
		}
		if i < len(arr)-1 {
			builder.WriteString(delim)
		}
	}
	return builder.String(), nil
}

// isOrderedMap checks if the argument is an instance of ordered map
func isOrderedMap(arg any) bool {
	switch arg.(type) {
	case orderedmap.OrderedMap[any, any], *orderedmap.OrderedMap[any, any]:
		return true
	default:
		return false
	}
}

// isCollectionType checks if the argument is either an array, a slice, a map or ordered map
func isCollectionType(arg any) bool {
	if arg == nil {
		return false
	}

	argType := reflect.TypeOf(arg).Kind()
	return isOrderedMap(arg) || argType == reflect.Slice || argType == reflect.Array || argType == reflect.Map
}

// aggregateValues extracts the stored value from different types of source:
// ordered map, map, slice and array. It gathers there values into an arr and returns it.
func aggregateValues(source any) []any {
	if isOrderedMap(source) {
		var om *orderedmap.OrderedMap[any, any]

		switch tmp := source.(type) {
		case orderedmap.OrderedMap[any, any]:
			// If source is an OrderedMap struct, use address of source
			om = &tmp
		case *orderedmap.OrderedMap[any, any]:
			om = tmp
		}

		arr := make([]any, 0, om.Len())
		for el := om.Front(); el != nil; el = el.Next() {
			arr = append(arr, el.Value)
		}
		return arr
	} else {
		v := reflect.ValueOf(source)
		arr := make([]any, 0, v.Len())

		switch v.Kind() {
		case reflect.Map:
			for _, value := range v.MapKeys() {
				arr = append(arr, v.MapIndex(value).Interface())
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				arr = append(arr, v.Index(i).Interface())
			}
		}
		return arr
	}
}
