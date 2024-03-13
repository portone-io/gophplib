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

type Floats interface {
	float32 | float64
}

// floatToString converts a float32 or float64 to a string based on the PHP 5.6 rules.
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
// NOTE : Converting float32 to float64 may lead to precision loss, hence using float64
// is recommended for higher accuracy.
// Reference :
//   - https://github.com/php/php-src/blob/php-5.6.40/Zend/zend_operators.c#L627-L633
func floatToString[T Floats](value T) string {
	f64 := float64(value)

	if math.IsNaN(f64) {
		return "NAN"
	}
	if math.IsInf(f64, 1) {
		return "INF"
	}
	if math.IsInf(f64, -1) {
		return "-INF"
	}
	return fmt.Sprintf("%.*G", 14, f64)
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

	// handle basic and composite types dynamically
	switch v := value.(type) {
	case string:
		return v, nil
	case bool:
		if v {
			return "1", nil
		} else {
			return "", nil
		}
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return floatToString(v), nil
	case float64:
		return floatToString(v), nil
	// check for special types such as a pointer of file, network, database resources
	case *os.File, *net.Conn, *sql.DB:
		// using a resource's address as the resource ID
		return fmt.Sprintf("Resource id %p", v), nil
	}
	// use reflection to handle array, slice, map types
	t := reflect.ValueOf(value).Kind()
	if t == reflect.Array || t == reflect.Slice || t == reflect.Map {
		return "Array", nil
	}
	if t == reflect.Struct {
		if result, ok := value.(toStringAble); ok {
			return result.toString(), nil
		}
	}
	// return an error for unsupported types.
	return "", fmt.Errorf("unsupported type : %T", value)
}
