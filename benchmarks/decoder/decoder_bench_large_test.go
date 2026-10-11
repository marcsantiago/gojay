package benchmarks

import (
	"testing"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/marcsantiago/gojay"
	"github.com/marcsantiago/gojay/benchmarks"
)

func BenchmarkJsonParserDecodeObjLarge(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = jsonparser.ArrayEach(benchmarks.LargeFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, _, _, _ = jsonparser.Get(value, "username")
			nothing()
		}, "users")

		_, _ = jsonparser.ArrayEach(benchmarks.LargeFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, _ = jsonparser.GetInt(value, "id")
			_, _, _, _ = jsonparser.Get(value, "slug")
			nothing()
		}, "topics", "topics")
	}
}

func BenchmarkJsonIterDecodeObjLarge(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.LargePayload{}
		_ = jsoniter.Unmarshal(benchmarks.LargeFixture, &result)
	}
}

func BenchmarkEasyJsonDecodeObjLarge(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.LargePayload{}
		_ = easyjson.Unmarshal(benchmarks.LargeFixture, &result)
	}
}

func BenchmarkGoJayDecodeObjLarge(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.LargePayload{}
		_ = gojay.UnmarshalJSONObject(benchmarks.LargeFixture, &result)
	}
}

func BenchmarkGoJayUnsafeDecodeObjLarge(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.LargePayload{}
		_ = gojay.Unsafe.UnmarshalJSONObject(benchmarks.LargeFixture, &result)
	}
}
