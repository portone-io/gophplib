package gophplib

import (
	"fmt"
	"reflect"
	"strings"
)

// Trim is a ported function that works exactly the same as PHP 5.6's trim
// function. For more information, see the [official PHP documentation].
//
// In PHP 5.6, when attempting to use the trim() function with a data type other
// than a string, it automatically converts the requested variable into a string
// before performing the trim. To achieve the same behavior in Go, this function
// converts the requested data types into strings and then utilize the trim
// function from the package strings. For more detailed information about the
// trim function in the package strings, see the [strings's trim documentation]
//
// NOTE: This function does not support the second parameter of original parse_str yet.
// It only strips the default characters (" \n\r\t\v\x00")
//
// References:
//   - https://www.php.net/manual/en/function.trim.php
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/string.c#L840-L850
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L425-L470
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L593-L661
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L261-L302
//
// Test Cases:
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/trim1.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/trim.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/trim_basic.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/strings/trim_variation1.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.trim.php
// [strings's trim documentation]: https://pkg.go.dev/strings#Trim
func Trim(value any) (ret string, err error) {
	// Convert a value to string
	characterString, err := zendParseArgAsString(value)
	if err != nil {
		err = fmt.Errorf("unsupported type : %s", reflect.TypeOf(value))
		return
	}

	const charSet = " \n\r\t\v\x00"

	ret = strings.Trim(characterString, charSet)
	return
}
