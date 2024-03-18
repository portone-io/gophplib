package gophplib

import (
	"database/sql"
	"fmt"
	"math"
	"net"
	"os"
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

// ConvertToString attempts to convert the given value to string, emulating PHP 5.6'S _convert_to_string behavior.
// Unlike PHP, which has built-in support for managing resource IDs for types like files and database connections,
// Go does not inherently manage resource IDs. Due to this language difference, this function uses the values' pointer
// address as the pseudo resource ID for identifiable resource types.
//
// This function returns error if given argument is not one of following:
// string, int, int8, int16, int32, int64, float32, float64, bool, nil, *os.File, *net.Conn, and *sql.DB,
// array, slice, map and any type which does not implement interface { toString() string }.
//
// Reference:
//   - _convert_to_string implementation:
//     https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L593-L661
//   - convert_object_to_type implementation:
//     https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L333-L357
func ConvertToString(value any) (string, error) {
	if value == nil {
		return "", nil
	}

	// use reflection to handle basic and composite types dynamically
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		// convert boolean to "1" or "" for true and false
		if v.Bool() {
			return "1", nil
		} else {
			return "", nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil
	case reflect.Float32, reflect.Float64:
		return floatToString(v.Float()), nil
	case reflect.Array, reflect.Slice, reflect.Map:
		return "Array", nil
	}

	// check for special types such as a pointer of file, network, database resources and
	// type which does not implement interface { toString() string }
	switch v := value.(type) {
	case toStringAble:
		return v.toString(), nil
	case *os.File, *net.Conn, *sql.DB:
		// using a resource's address as the resource ID
		return fmt.Sprintf("Resource id %p", v), nil
	default:
		// return an error for unsupported types.
		return "", fmt.Errorf("unsupported type : %T", value)
	}
}
