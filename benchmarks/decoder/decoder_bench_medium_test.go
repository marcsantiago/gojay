package benchmarks

import (
	"encoding/json"
	"testing"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/marcsantiago/gojay"
	"github.com/marcsantiago/gojay/benchmarks"
)

func BenchmarkJsonIterDecodeObjMedium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.MediumPayload{}
		_ = jsoniter.Unmarshal(benchmarks.MediumFixture, &result)
	}
}

/*
github.com/buger/jsonparser
*/
func BenchmarkJSONParserDecodeObjMedium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _, _ = jsonparser.Get(benchmarks.MediumFixture, "person", "name", "fullName")
		_, _ = jsonparser.GetInt(benchmarks.MediumFixture, "person", "github", "followers")
		_, _, _, _ = jsonparser.Get(benchmarks.MediumFixture, "company")

		_, _ = jsonparser.ArrayEach(benchmarks.MediumFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, _, _, _ = jsonparser.Get(value, "url")
			nothing()
		}, "person", "gravatar", "avatars")
	}
}

func BenchmarkEncodingJsonStructMedium(b *testing.B) {
	for b.Loop() {
		data := benchmarks.MediumPayload{}
		_ = json.Unmarshal(benchmarks.MediumFixture, &data)
	}
}

func BenchmarkEasyJsonDecodeObjMedium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.MediumPayload{}
		_ = easyjson.Unmarshal(benchmarks.MediumFixture, &result)
	}
}

func BenchmarkGoJayDecodeObjMedium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.MediumPayload{}
		err := gojay.UnmarshalJSONObject(benchmarks.MediumFixture, &result)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGoJayUnsafeDecodeObjMedium(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.MediumPayload{}
		err := gojay.Unsafe.UnmarshalJSONObject(benchmarks.MediumFixture, &result)
		if err != nil {
			b.Error(err)
		}
	}
}
