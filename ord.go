package gophplib

import (
	"fmt"
	"reflect"
)

// Ord is a ported functions that works exactly the same as PHP 5.6's ord function.
// In PHP 5.6, when the ord() function is used with a data type other
// than a string, it automatically converts the given variable into a string
// before processing it. To achieve the same behavior in Go,
// this function converts an argument to string using the zendParseArgAsString() function.
// For more information, see the [official PHP documentation].
//
// This function returns error if given argument is not one of following:
// string, int, int64, float64, bool, nil, and any type which does not implement
// interface { toString() string }.
//
// Reference :
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/string.c#L2666-L2676
//
// Test Cases:
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/ord_basic.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/ord_error.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/ord_variation1.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.ord.php
func Ord(character any) (byte, error) {
	// Convert a character to string
	characterString, err := zendParseArgAsString(character)
	if err != nil {
		return 0, fmt.Errorf("unsupported type : %s", reflect.TypeOf(character))
	}

	// Check if the characterString is not empty
	if len(characterString) > 0 {
		// Return the first byte of input argument's string representation
		return []byte(characterString)[0], nil
	}
	// Return for empty strings
	return 0, nil
}
