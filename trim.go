package gophplib

import (
	"fmt"
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
func Trim(value any) any {
	var charSet = " \n\r\t\v\x00"

	// Handle different types of value
	switch v := value.(type) {
	case string:
		return strings.Trim(v, charSet)
	case int:
		return strings.Trim(fmt.Sprintf("%v", value), charSet)
	case float64:
		// Rules for converting a float to a string in PHP 5.6:
		// - Allows up to a maximum of 14 digits, including both integer and decimal places.
		// - Remove trailing zeros from the fractional part
		//	 ex) 123.4000 → "123.4"
		// - Keep the values as is if the last digit is not 0.
		//   ex) 123.45 → "123.45"
		// - If the integer part exceeds 14 digits, use exponential notation.
		//   ex) 123456789123456.40 → "1.2345678901234e+14"
		// - If the total number of digits exceeds 14, truncate the decimal places.
		//   ex) 123.45678901234 → "123.4567890123"
		result := fmt.Sprintf("%.*G", 14, v)
		return strings.Trim(result, charSet)
	case bool:
		// return "1" for true and an empty string("") for false
		if v {
			return "1"
		} else {
			return ""
		}
	case nil:
		return ""
	case toStringAble:
		// For types implementing toString(), trim the return value of toString()
		return strings.Trim(v.toString(), charSet)
	default:
		// For other types, return nil
		return nil
	}
}
