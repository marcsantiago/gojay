package gojay

import (
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	invalidCodePoint = InvalidJSONError("Invalid unicode code point")
	invalidJson      = InvalidJSONError("Invalid JSON")

	unicodeReplacement = make([]byte, 3)
	_                  = utf8.EncodeRune(unicodeReplacement, unicode.ReplacementChar)
)

func (dec *Decoder) getUnicode() (rune, error) {
	i := 0
	r := rune(0)
	for ; (dec.cursor < dec.length || dec.read()) && i < 4; dec.cursor++ {
		c := dec.data[dec.cursor]
		if c >= '0' && c <= '9' {
			r = r*16 + rune(c-'0')
		} else if c >= 'a' && c <= 'f' {
			r = r*16 + rune(c-'a'+10)
		} else if c >= 'A' && c <= 'F' {
			r = r*16 + rune(c-'A'+10)
		} else {
			return 0, invalidCodePoint
		}
		i++
	}
	return r, nil
}

func (dec *Decoder) appendEscapeChar(str []byte, c byte) ([]byte, error) {
	switch c {
	case 't':
		str = append(str, '\t')
	case 'n':
		str = append(str, '\n')
	case 'r':
		str = append(str, '\r')
	case 'b':
		str = append(str, '\b')
	case 'f':
		str = append(str, '\f')
	case '\\':
		str = append(str, '\\')
	default:
		return nil, invalidJson
	}
	return str, nil
}

func (dec *Decoder) parseUnicode() ([]byte, error) {
	// get unicode after u
	r, err := dec.getUnicode()
	if err != nil {
		if err == invalidCodePoint {
			if dec.disableStrictUnicode {
				return unicodeReplacement, nil
			}
		}
		return nil, err
	}
	// no error start making new string
	str := make([]byte, 16, 16)
	i := 0
	// check if code can be a surrogate utf16
	if utf16.IsSurrogate(r) {
		if dec.cursor >= dec.length && !dec.read() {
			return nil, dec.raiseInvalidJSONErr(dec.cursor)
		}
		c := dec.data[dec.cursor]
		if c != '\\' {
			i += utf8.EncodeRune(str, r)
			return str[:i], nil
		}
		dec.cursor++
		if dec.cursor >= dec.length && !dec.read() {
			return nil, dec.raiseInvalidJSONErr(dec.cursor)
		}
		c = dec.data[dec.cursor]
		if c != 'u' {
			i += utf8.EncodeRune(str, r)
			str, err = dec.appendEscapeChar(str[:i], c)
			if err != nil {
				dec.err = err
				return nil, err
			}
			i++
			dec.cursor++
			return str[:i], nil
		}
		dec.cursor++
		r2, err := dec.getUnicode()
		if err != nil {
			return nil, err
		}
		combined := utf16.DecodeRune(r, r2)
		if combined == '\uFFFD' {
			i += utf8.EncodeRune(str, r)
			i += utf8.EncodeRune(str, r2)
		} else {
			i += utf8.EncodeRune(str, combined)
		}
		return str[:i], nil
	}
	i += utf8.EncodeRune(str, r)
	return str[:i], nil
}
