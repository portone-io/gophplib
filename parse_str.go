package gophplib

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

// ParseStr is a ported function that works exactly the same as PHP's parse_str
// function. For more information, see the [official PHP documentation].
//
// All keys of the returned map are either string or int. Type of all values of
// the returned map (RetVal) are either string or "map[string | int]RetVal".
//
// Reference:
//   - https://www.php.net/manual/en/function.parse-str.php
//   - https://github.com/php/php-src/blob/php-5.6.40/main/php_variables.c#L450-L496
//   - https://github.com/php/php-src/blob/php-8.3.0/main/php_variables.c#L523-L568
//
// [official PHP documentation]: https://www.php.net/manual/en/function.parse-str.php
func ParseStr(input string) map[any]any {
	ret := newPhpArray()

	// Split input with '&'
	pairs := strings.Split(input, "&")

	for _, pair := range pairs {
		// Skip empty pair
		if pair == "" {
			continue
		}

		// Cut pair with '='
		key, value, _ := strings.Cut(pair, "=")
		registerVariableSafe(Urldecode(key), Urldecode(value), ret)
	}

	return ret.intoMap()
}

// registerVariableSafe is a ported function that works exactly the same as
// PHP's php_register_variable_safe function.
//
// Reference:
//   - https://github.com/php/php-src/blob/php-5.6.40/main/php_variables.c#L59-L233
//   - https://github.com/php/php-src/blob/php-8.3.0/main/php_variables.c#L90-L314
//
// TODO: Add tests
func registerVariableSafe(key, value string, track *phpArray) {
	// NOTE: key is "var_name", value is "val", track is "track_vars_array" in
	// below PHP verion's function signature.
	//
	// PHPAPI void php_register_variable_ex(const char *var_name, zval *val, zval *track_vars_array)

	// ignore leading spaces in the variable name
	strings.TrimLeft(key, " ")

	// Prepare variable name
	// NOTE: key_new is "var" and "var_orig" in the original PHP codes.
	key_new := []byte(key)

	// ensure that we don't have spaces or dots in the variable name (not binary safe)
	is_array := false
	index_slice := []byte(nil) // index_slice is "ip" in the original PHP codes.
	for i, c := range key_new {
		if c == ' ' || c == '.' {
			key_new[i] = '_'
		} else if c == '[' {
			is_array = true
			key_new, index_slice = key_new[:i], key_new[i:]
			break
		}
	}

	// empty variable name, or variable name with a space in it
	if len(key_new) == 0 {
		return
	}

	index := key_new

	if is_array {
		// We do not perform max nesting level check here
		idx := 0 // idx is offset of "ip" pointer in the original PHP codes.
		for true {
			idx++
			idx_s := idx // idx_next is "index_s" in the original PHP codes.
			if isAsciiSpace(index_slice[idx]) {
				idx++
			}
			if index_slice[idx] == ']' {
				idx_s = -1
			} else {
				ret := bytes.IndexByte(index_slice[idx:], ']')
				if ret == -1 {
					// not an index; un-terminate the var name
					index_slice[idx_s-1] = '_'

					// NOTE: concat of `index` and `index_slice[idx_s-1:]` only
					// occurs when idx_s == 1.
					if index != nil && idx_s == 1 {
						index = append(index, index_slice...)
					}

					goto plain_var
				}
				idx += ret
			}

			var subdict *phpArray
			if index == nil {
				subdict = newPhpArray()
				track.setNext(subdict)
			} else {
				value, ok := track.get(index)
				if !ok {
					subdict = newPhpArray()
					track.set(index, subdict)
				} else {
					// References for origianl PHP codes of here:
					//   - https://www.phpinternalsbook.com/php7/zvals/memory_management.html
					//   - https://www.phpinternalsbook.com/php7/zvals/basic_structure.html
					underlying, ok := value.(*phpArray)
					if !ok {
						subdict = newPhpArray()
						track.set(index, subdict)
					} else {
						subdict = underlying
					}
				}
			}
			track = subdict
			if idx_s != -1 {
				index = index_slice[idx_s:idx]
			} else {
				index = nil
			}

			idx++
			if idx < len(index_slice) && index_slice[idx] == '[' {
				// Do nothing
			} else {
				goto plain_var
			}
		}
	}
plain_var:
	if index == nil {
		track.setNext(value)
	} else {
		track.set(index, value)
	}
}

// phpArray is a map[any]any which behaves like PHP's array. It maintains
// internal next (i.e. nNextFreeElement of PHP) state and it automatically
// converts numeric string keys to integer keys.
type phpArray struct {
	next int
	// Key is either string or int.
	// Value is either string or *phpArray.
	d map[any]any
}

func newPhpArray() *phpArray {
	return &phpArray{
		next: 0,
		d:    make(map[any]any),
	}
}

// It returns a map[any]any whose keys are either string or int, and whose
// values (RetVal) are either string or "map[string | int]RetVal".
//
// TODO: Add tests
func (p *phpArray) intoMap() map[any]any {
	ret := make(map[any]any)
	for k, v := range p.d {
		if sub, ok := v.(*phpArray); ok {
			ret[k] = sub.intoMap()
		} else {
			ret[k] = v
		}
	}
	return ret
}

// TODO: Add tests
func (p *phpArray) get(key []byte) (any, bool) {
	k := zendHandleNumericStr(key)
	v, ok := p.d[k]
	return v, ok
}

// TODO: Add tests
func (p *phpArray) set(key []byte, value any) {
	k := zendHandleNumericStr(key)
	if numeric, ok := k.(int); ok {
		p.next = maxInt(p.next, numeric+1)
	}
	p.d[k] = value
}

// TODO: Add tests
func (p *phpArray) setNext(value any) {
	p.d[p.next] = value
	p.next++
}

// TODO: Add tests
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// References:
//   - https://github.com/php/php-src/blob/php-8.3.0/Zend/zend_hash.h#L388-L404
//   - https://github.com/php/php-src/blob/php-8.3.0/Zend/zend_hash.c#L3262-L3299
//
// TODO: Add tests
func zendHandleNumericStr(input []byte) any {
	str := string(input)
	// Test if input is numeric string without leading zeros or a plus sign
	if !regexp.MustCompile(`^-?[1-9][0-9]*$|^0$`).MatchString(str) {
		return str
	}
	num, err := strconv.Atoi(str)
	// Check if it overflows
	if err != nil {
		return str
	}
	return num
}

// isAsciiSpace is an ASCII-only version of C's isspace.
//
// References:
//   - https://en.cppreference.com/w/c/string/byte/isspace
//
// TODO: Add tests
func isAsciiSpace(c byte) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}
