package gojay

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeNull(t *testing.T) {
	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st element",
			baseJSON:     `[`,
			expectedJSON: `[null,null`,
		},
		{
			name:         "basic last element",
			baseJSON:     `["test"`,
			expectedJSON: `["test",null,null`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.Null()
			enc.AddNull()
			_, _ = enc.Write()
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncodeNullKey(t *testing.T) {
	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st element",
			baseJSON:     `{`,
			expectedJSON: `{"foo":null,"bar":null`,
		},
		{
			name:         "basic last element",
			baseJSON:     `{"test":"test"`,
			expectedJSON: `{"test":"test","foo":null,"bar":null`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.NullKey("foo")
			enc.AddNullKey("bar")
			_, _ = enc.Write()
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}
