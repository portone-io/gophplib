package gophplib

import (
	"fmt"
	"math"
	"reflect"
)

type toStringAble interface {
	toString() string
}

// floatToString converts a float64 to a string based on the PHP 5.6 rules.
//   - Allows up to a maximum of 14 digits, including both integer and decimal places.
//   - Remove trailing zeros from the fractional part
//     ex) 123.4000 → "123.4"
//   - Keep the values as is if the last digit is not 0.
//     ex) 123.45 → "123.45"
//   - If the integer part exceeds 14 digits, use exponential notation.
//     ex) 123456789123456.40 → "1.2345678901234e+14"
//   - If the total number of digits exceeds 14, truncate the decimal places.
//     ex) 123.45678901234 → "123.4567890123"
//
// Reference :
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L627-L633
func floatToString(value float64) string {
	if math.IsNaN(value) {
		return "NAN"
	}
	if math.IsInf(value, 1) {
		return "INF"
	}
	if math.IsInf(value, -1) {
		return "-INF"
	}
	return fmt.Sprintf("%.*G", 14, value)
}

// ZendParseArgImpl attempts to replicate the behavior of the 'zend_parse_arg_impl' function
// from PHP 5.6, specifically for the case where the 'spec' parameter is "s".
// It handles conversion of different types to string in a way that aligns with PHP's type juggling rules.
//
// Reference :
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L685-L713
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L425-L470
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L593-L661
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L261-L301
func ZendParseArgImpl(value any) (string, error) {
	var str string

	switch v := value.(type) {
	case string:
		str = v
	case int, int64:
		str = fmt.Sprintf("%v", v)
	case float64:
		str = floatToString(v)
	case bool:
		// return "1" for true and an empty string("") for false
		if v {
			str = "1"
		} else {
			str = ""
		}
	case nil:
		// TODO: handle check_null
		str = ""
	case toStringAble:
		// For types implementing toString(), get the value of toString()
		str = v.toString()
	default:
		return "", fmt.Errorf("unsupported type : %s", reflect.TypeOf(v))
	}
	return str, nil
}
