package gojay

import (
	"unsafe"
)

// DecodeString reads the next JSON-encoded value from the decoder's input (io.Reader) and stores it in the string pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON into a Go value.
func (dec *Decoder) DecodeString(v *string) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeString(v)
}

func (dec *Decoder) decodeString(v *string) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '"':
			dec.cursor++
			start, end, err := dec.getString()
			if err != nil {
				return err
			}
			d := dec.data[start : end-1]
			*v = *(*string)(unsafe.Pointer(&d))
			dec.cursor = end
			return nil
		case 'n':
			dec.cursor++
			if err := dec.assertNull(); err != nil {
				return err
			}
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			if err := dec.skipData(); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (dec *Decoder) decodeStringNull(v **string) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '"':
			dec.cursor++
			start, end, err := dec.getString()
			if err != nil {
				return err
			}
			if *v == nil {
				*v = new(string)
			}
			d := dec.data[start : end-1]
			**v = *(*string)(unsafe.Pointer(&d))
			dec.cursor = end
			return nil
		case 'n':
			dec.cursor++
			if err := dec.assertNull(); err != nil {
				return err
			}
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			if err := dec.skipData(); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

// getString extracts a JSON string and returns its decoded-content start and
// the position immediately after its closing quote.
//
// Escaped characters are decoded in place using separate read and write
// cursors. The trailing JSON data is moved at most once per escaped string.
func (dec *Decoder) getString() (int, int, error) {
	keyStart := dec.cursor
	read := dec.cursor
	write := dec.cursor

	for read < dec.length || dec.read() {
		switch dec.data[read] {
		case '"':
			if write != read {
				remaining := dec.length - read
				copy(dec.data[write:write+remaining], dec.data[read:dec.length])
				dec.length -= read - write
			}

			dec.cursor = write + 1
			return keyStart, dec.cursor, nil

		case '\\':
			// Advance to the escaped character.
			read++
			if read >= dec.length && !dec.read() {
				dec.cursor = read
				return 0, 0, dec.raiseInvalidJSONErr(dec.cursor)
			}

			switch dec.data[read] {
			case '"':
				dec.data[write] = '"'
				write++
			case '\\':
				dec.data[write] = '\\'
				write++
			case '/':
				dec.data[write] = '/'
				write++
			case 'b':
				dec.data[write] = '\b'
				write++
			case 'f':
				dec.data[write] = '\f'
				write++
			case 'n':
				dec.data[write] = '\n'
				write++
			case 'r':
				dec.data[write] = '\r'
				write++
			case 't':
				dec.data[write] = '\t'
				write++
			case 'u':
				// parseUnicode expects the cursor at the first hex digit.
				dec.cursor = read + 1
				decoded, err := dec.parseUnicode()
				if err != nil {
					return 0, 0, err
				}

				// parseUnicode leaves the cursor immediately after the
				// Unicode escape (including a surrogate pair, if present).
				read = dec.cursor
				write += copy(dec.data[write:], decoded)
				continue
			default:
				dec.cursor = read
				return 0, 0, dec.raiseInvalidJSONErr(dec.cursor)
			}

			// Consume the escaped character exactly once. read now points
			// at the next original input byte.
			read++

		default:
			if write != read {
				dec.data[write] = dec.data[read]
			}
			write++
			read++
		}
	}

	dec.cursor = read
	return 0, 0, dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) skipEscapedString() error {
	start := dec.cursor
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		if dec.data[dec.cursor] != '\\' {
			d := dec.data[dec.cursor]
			dec.cursor++
			nSlash := dec.cursor - start
			switch d {
			case '"':
				if nSlash&1 != 1 {
					return dec.raiseInvalidJSONErr(dec.cursor)
				}
				return nil
			case 'u':
				if err := dec.skipString(); err != nil {
					return err
				}
				dec.cursor--
				return nil
			case 'n', 'r', 't', '/', 'f', 'b':
				return nil
			default:
				if nSlash&1 == 1 {
					return dec.raiseInvalidJSONErr(dec.cursor)
				}
				return nil
			}
		}
	}
	return dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) skipString() error {
	for dec.cursor < dec.length || dec.read() {
		switch dec.data[dec.cursor] {
		case '"':
			dec.cursor++
			return nil
		case '\\':
			dec.cursor++
			if err := dec.skipEscapedString(); err != nil {
				return err
			}
		default:
			dec.cursor++
		}
	}
	return dec.raiseInvalidJSONErr(len(dec.data) - 1)
}

// Add Values functions

// AddString decodes the JSON value within an object or an array to a *string.
// If next key is not a JSON string nor null, InvalidUnmarshalError will be returned.
func (dec *Decoder) AddString(v *string) error {
	return dec.String(v)
}

// AddStringNull decodes the JSON value within an object or an array to a *string.
// If next key is not a JSON string nor null, InvalidUnmarshalError will be returned.
// If a `null` is encountered, gojay does not change the value of the pointer.
func (dec *Decoder) AddStringNull(v **string) error {
	return dec.StringNull(v)
}

// String decodes the JSON value within an object or an array to a *string.
func (dec *Decoder) String(v *string) error {
	err := dec.decodeString(v)
	if err != nil {
		return err
	}
	dec.called |= 1
	return nil
}

// StringNull decodes the JSON value within an object or an array to a **string.
func (dec *Decoder) StringNull(v **string) error {
	err := dec.decodeStringNull(v)
	if err != nil {
		return err
	}
	dec.called |= 1
	return nil
}
