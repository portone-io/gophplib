package gophplib

import (
	"fmt"
	"reflect"
	"strings"
)

// Implode replicates the behavior of PHP 5.6's implode function in GO.
// This function concatenates the elements of an array into a single string using a specified separator.
// For more information, see the [official PHP documentation].
//
// The function supports flexible argument patterns to accommodate various use cases.
//   - arg1: Can be either a string (used as the separator) or an array of elements to be joined.
//   - arg2: Optional. When provided, it expects the first element to serve as the array of elements to be joined
//     if arg1 is a string, or as the separator if arg1 is an array.
//
// Behavior:
//
// 1. If only arg1 is provided:
//   - If arg1 is an array, the function joins the elements using an empty string as the default separator.
//   - If arg1 is not an array, the function returns an error.
//
// 2. If both arg1 and arg2 are provided:
//   - If arg1 is an array, the first element of arg2 is converted to string and used as the separator.
//   - If arg1 is not an array and the first element of arg2 is an array, arg1 is converted to a string
//     and used as the separator, with the first of element of arg2 being the array to implode.
//   - If neither arg1 nor the first of element of arg2 is an array, the function returns an error.
//
// Non-string elements within the array are converted to strings using a ConvertToString function
// before joining.
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
func Implode(arg1 any, arg2 ...any) (string, error) {
	var delim string
	var arr []any

	var IsArrayOrSlice = func(arg any) bool {
		if arg != nil {
			argType := reflect.TypeOf(arg).Kind()
			return argType == reflect.Slice || argType == reflect.Array
		}
		return false
	}

	IsArg1ArrayOrSlice := IsArrayOrSlice(arg1)
	// Check if arg2 is not provided
	if len(arg2) == 0 {
		if IsArg1ArrayOrSlice {
			v := reflect.ValueOf(arg1)
			for i := 0; i < v.Len(); i++ {
				arr = append(arr, v.Index(i).Interface())
			}
			delim = ""
		} else {
			return "", fmt.Errorf("argument must be an array, but got %v", reflect.TypeOf(arg1))
		}
	} else {
		IsArg2ArrayOrSlice := IsArrayOrSlice(arg2[0])
		if IsArg1ArrayOrSlice {
			delim, _ = ConvertToString(arg2[0])
			v := reflect.ValueOf(arg1)
			for i := 0; i < v.Len(); i++ {
				arr = append(arr, v.Index(i).Interface())
			}
		} else if !IsArg1ArrayOrSlice && IsArg2ArrayOrSlice {
			delim, _ = ConvertToString(arg1)
			v := reflect.ValueOf(arg2[0])
			for i := 0; i < v.Len(); i++ {
				arr = append(arr, v.Index(i).Interface())
			}
		} else {
			return "", fmt.Errorf("invalid arguments passed, got %v, %v", reflect.TypeOf(arg1), reflect.TypeOf(arg2[0]))
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
