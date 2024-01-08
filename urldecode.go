package gophplib

import (
	"strings"
)

// Urldecode is a ported function that works exactly the same as PHP's urldecode
// function. For more information, see the [official PHP documentation].
//
// Unlike net/url's QueryUnescape function, this function *never* fails. Instead
// of returning an error, it leaves invalid percent encoded sequences as is.
// And contrary to its name, it does not follow percent encoding specification
// of RFC 3986 since it decodes '+' to ' '. This is done to be compatible with
// PHP's urldecode function.
//
// References:
//   - https://www.php.net/manual/en/function.urldecode.php
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/url.c#L513-L561
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/url.c#L578-L618
//
// [official PHP documentation]: https://www.php.net/manual/en/function.urldecode.php
func Urldecode(input string) string {
	buf := []byte(input)
	length := len(buf)

	j := 0
	for i := 0; i < length; i, j = i+1, j+1 {
		if buf[i] == '+' {
			buf[j] = ' '
		} else if buf[i] == '%' && i+2 < length && isxdigit(buf[i+1]) && isxdigit(buf[i+2]) {
			buf[j] = htoi(buf[i+1], buf[i+2])
			i += 2
		} else {
			buf[j] = buf[i]
		}
	}

	return strings.ToValidUTF8(string(buf[:j]), "�")
}

// isxdigit is a ported function that works exactly the same as C's isxdigit
// function.
//
// References:
//   - https://en.cppreference.com/w/c/string/byte/isxdigit
func isxdigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// htoi is a ported function that works exactly the same as PHP's php_htoi
// function. It returns the rune value of the hexadecimal number represented by
// the two runes hi and lo.
//
// It expects both hi and lo to be valid hexadecimal digits. (ex: '0'-'9',
// 'a'-'f', 'A'-'F') Otherwise, the result is undefined.
//
// References:
//   - https://github.com/php/php-src/blob/php-8.3.0/ext/standard/url.c#L426-L444
//   - https://github.com/php/php-src/blob/php-5.6.40/ext/standard/url.c#L407-L426
func htoi(hi, lo byte) byte {
	if 'A' <= hi && hi <= 'Z' {
		hi += 'a' - 'A'
	}
	if '0' <= hi && hi <= '9' {
		hi -= '0'
	} else {
		hi -= 'a' - 10
	}

	if 'A' <= lo && lo <= 'Z' {
		lo += 'a' - 'A'
	}
	if '0' <= lo && lo <= '9' {
		lo -= '0'
	} else {
		lo -= 'a' - 10
	}

	return hi*16 + lo
}
