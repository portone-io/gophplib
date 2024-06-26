package gophplib

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

// ParseStr is a ported function that works exactly the same as PHP's parse_str
// function. For more information, see the [official PHP documentation].
//
// All keys of the returned orderedmap.OrderedMap are either string or int. Type of all values of
// the returned orderedmap.OrderedMap are either string or "orderedmap.OrderedMap[string | int]RetVal".
// For more information about orderedmap libaray, see the [orderedmap documentation].
//
// Reference:
//   - https://www.php.net/manual/en/function.parse-str.php
//   - https://github.com/php/php-src/blob/php-5.6.40/main/php_variables.c#L450-L496
//   - https://github.com/php/php-src/blob/php-8.3.0/main/php_variables.c#L523-L568
//
// [official PHP documentation]: https://www.php.net/manual/en/function.parse-str.php
// [orderedmap documentation]: https://pkg.go.dev/github.com/elliotchance/orderedmap/v2@v2.2.0
func ParseStr(input string) orderedmap.OrderedMap[any, any] {
	ret := newPHPArray()

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
func registerVariableSafe(key, value string, track *phpSymtable) {
	// NOTE: key is "var_name", value is "val", track is "track_vars_array" in
	// below PHP version's function signature.
	//
	// PHPAPI void php_register_variable_ex(const char *var_name, zval *val, zval *track_vars_array)

	// ignore leading spaces in the variable name
	key = strings.TrimLeft(key, " ")

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
		for {
			idx++
			idx_s := idx // idx_next is "index_s" in the original PHP codes.
			if isAsciiWhitespace(index_slice[idx]) {
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

			var subdict *phpSymtable
			if index == nil {
				subdict = newPHPArray()
				track.setNext(subdict)
			} else {
				value, ok := track.get(index)
				if !ok {
					subdict = newPHPArray()
					track.set(index, subdict)
				} else {
					// References for origianl PHP codes of here:
					//   - https://www.phpinternalsbook.com/php7/zvals/memory_management.html
					//   - https://www.phpinternalsbook.com/php7/zvals/basic_structure.html
					underlying, ok := value.(*phpSymtable)
					if !ok {
						subdict = newPHPArray()
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

// phpSymtable is a orderedmap.OrderedMap[any, any] which behaves like PHP's array. It maintains
// internal next (i.e. nNextFreeElement of PHP) state and it automatically
// converts numeric string keys to integer keys.
type phpSymtable struct {
	next int
	// Key is either string or int.
	// Value is either string or *phpSymtable.
	d orderedmap.OrderedMap[any, any]
}

func newPHPArray() *phpSymtable {
	return &phpSymtable{
		next: 0,
		d:    *orderedmap.NewOrderedMap[any, any](),
	}
}

// It returns an orderedmap.OrderedMap[any, any] whose keys are either string or int, and whose
// values (RetVal) are either string or "orderedmap.OrderedMap[string | int, RetVal]".
func (p *phpSymtable) intoMap() orderedmap.OrderedMap[any, any] {
	ret := *orderedmap.NewOrderedMap[any, any]()
	for el := p.d.Front(); el != nil; el = el.Next() {
		key := el.Key
		value := el.Value
		if sub, ok := el.Value.(*phpSymtable); ok {
			ret.Set(key, sub.intoMap())
		} else {
			ret.Set(key, value)
		}
	}
	return ret
}

func (p *phpSymtable) get(key []byte) (any, bool) {
	k := phpNumericOrString(key)
	v, ok := p.d.Get(k)
	return v, ok
}

func (p *phpSymtable) set(key []byte, value any) {
	k := phpNumericOrString(key)
	if numeric, ok := k.(int); ok {
		p.next = maxInt(p.next, numeric+1)
	}
	p.d.Set(k, value)
}

func (p *phpSymtable) setNext(value any) {
	p.d.Set(p.next, value)
	p.next++
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func phpNumericOrString(input []byte) any {
	str := string(input)
	if !zendHandleNumericStr(str) {
		return str
	}
	num, err := strconv.Atoi(str)
	// Check if it overflows
	if err != nil {
		return str
	}
	return num
}

// zendHandleNumericStr is a ported function that works exactly the same as
// PHP's _zend_handle_numeric_str function.
//
// It returns true if the input string meets all the following conditions:
//   - It is a signed integer string without leading zeros. (positive sign is
//     not allowed, only negative sign is allowed)
//   - It is not a negative zero.
//
// It behaves same with regexp.MustCompile(`^-?[1-9][0-9]*$|^0$`).MatchString
//
// References:
//   - https://github.com/php/php-src/blob/php-8.3.0/Zend/zend_hash.h#L388-L404
//   - https://github.com/php/php-src/blob/php-8.3.0/Zend/zend_hash.c#L3262-L3299
func zendHandleNumericStr(s string) bool {
	// Handle few cases first to make further checks simpler
	switch s {
	case "0":
		return true
	case "", "-":
		return false
	}

	// Check for negative sign
	begin := 0
	if s[0] == '-' {
		begin = 1
	}

	// Ensure the first character isn't '0'
	if s[begin] == '0' {
		return false
	}

	// Check that all characters are digits
	for _, ch := range s[begin:] {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

// isAsciiWhitespace is an ASCII-only version of C's isspace.
//
// References:
//   - https://en.cppreference.com/w/c/string/byte/isspace
func isAsciiWhitespace(c byte) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}
