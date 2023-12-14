package gophplib

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleParseStr() {
	fmt.Println(ParseStr("key=value&foo=bar"))
	// Output: map[foo:bar key:value]
}

func TestParseStr(t *testing.T) {
	result := ParseStr("key=value&foo=bar")
	if !reflect.DeepEqual(result, map[string]interface{}{
		"key": "value",
		"foo": "bar",
	}) {
		t.Error("ParseStr() failed")
	}
}

/*

TODO:
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic1.phpt
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic2.phpt
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic3.phpt
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_basic4.phpt
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/parse_str_memory_error.phpt
https://github.com/php/php-src/blob/php-8.3.0/ext/standard/tests/strings/bug77439.phpt
https://github.com/simnalamburt/snippets/blob/3f095f671ca7277a9dabfe60e16fb749effb2e7c/php/parse_str.php

*/
