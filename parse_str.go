package gophplib

// ParseStr is a ported function that works exactly the same as PHP's parse_str
// function. For more information, see the [official PHP documentation].
//
// Reference:
//   - https://www.php.net/manual/en/function.parse-str.php
//   - https://github.com/php/php-src/blob/php-5.6.40/main/php_variables.c#L450-L496
//
// Test cases:
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic1.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic2.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic3.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic4.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_memory_error.phpt
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/bug77439.phpt
//
// [official PHP documentation]: https://www.php.net/manual/en/function.parse-str.php
func ParseStr(input string) map[interface{}]interface{} {
	// TODO: Implementation
	return map[interface{}]interface{}{
		"key": "value",
		"foo": "bar",
	}
}
