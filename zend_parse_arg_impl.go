package gophplib

import (
	"fmt"
	"reflect"
)

// zendParseArgAsString attempts to replicate the behavior of the 'zend_parse_arg_impl' function
// from PHP 5.6, specifically for the case where the 'spec' parameter is "s".
// It handles conversion of different types to string in a way that aligns with PHP's type juggling rules.
//
// This function returns error if given argument is not one of following:
// string, int, int8, int16, int32, int64, float32, float64, bool, nil
// and any type which does not implement interface { toString() string }.
//
// Reference :
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L685-L713
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L425-L470
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L593-L661
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_API.c#L261-L301
func zendParseArgAsString(value any) (string, error) {
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
