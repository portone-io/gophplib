package gophplib

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
)

// Base64Encode emulates the functionality of PHP 5.6's base64_encode function.
// For more information, see the [official PHP documentation].
// In PHP 5.6, an error is triggered if the input size is excessively large due to memory limitations.
// This Go implementation includes similar checks to emulate PHP's memory limitation conditions.
// Additionally, this function converts different types of variables to a string, following PHP's dynamic typing approach.
// After ensuring the memory constraints are met and converting the input to a string,
// it uses the EncodeToString function from the encoding/base64 package to perform the Base64 encoding.
// The result is a Base64 encoded string, consistent with PHP's output for the same input.
// For more detailed information about the EncodeToString function in the package encoding/base64,
// see the [encoding/base64's EncodeToString documentation]
//
// This function returns error if given argument is not one of following:
// string, int, int64, float64, bool, nil, and any type which does not implement
// interface { toString() string }.
//
// PHP references:
//   - base64_encode definition:
//     https://github.com/php/php-src/blob/php-5.6.40/ext/standard/base64.c#L224-L241
//   - base64_encode implementation:
//     https://github.com/php/php-src/blob/4b8f72da5dfb201af4e82dee960261d8657e414f/ext/standard/base64.c#L56-L106
//
// Test Cases :
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/url/base64_encode_basic_001.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/url/base64_encode_basic_002.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/url/base64_encode_error_001.phpt
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/tests/url/base64_encode_variation_001.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.base64-encode
// [encoding/base64's EncodeToString documentation]: https://pkg.go.dev/encoding/base64#Encoding.EncodeToString
func Base64Encode(value any) (string, error) {
	// Convert a value to string
	characterString, err := zendParseArgAsString(value)
	if err != nil {
		return "", fmt.Errorf("unsupported type : %s", reflect.TypeOf(value))
	}

	// Base64 encoding converts 3 bytes into 4 bytes ASCII characters,
	// and adds padding 2 bytes if the converted data is not a multiple of 3.
	// In PHP, the base64_encode function includes a memory allocation limit check
	// to prevent potential overflow issues due to the increase in data size during encoding.
	// In Go, such checks are generally not required thanks to its robust memory management system.
	// However, to maintain exact behavioral parity with the PHP implementation,
	// this function includes memory limit check too.
	if (len(characterString)+2)/3 > math.MaxInt32/4 {
		return "", fmt.Errorf("string too long, maximum is 1610612733")
	}
	encodedString := base64.StdEncoding.EncodeToString([]byte(characterString))
	return encodedString, nil
}
