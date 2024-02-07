package gophplib

import (
	"fmt"
	"reflect"
)

// Strlen is a ported function that works exactly the same as PHP 5.6's strlen function.
// In PHP 5.6, when the strlen() function is used with a data type other
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
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_builtin_functions.c#L479-L492
//
// Test Case :
//   - https://github.dev/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/strlen.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/strlen_variation1.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/strlen_error.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/strlen_basic.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.strlen.php
func Strlen(value any) (int, error) {
	// Convert a value to string
	characterString, err := zendParseArgAsString(value)
	if err != nil {
		return 0, fmt.Errorf("unsupported type : %s", reflect.TypeOf(value))
	}
	return len(characterString), nil
}
