package benchmarks

import (
	"encoding/json"
	_ "fmt"
	"testing"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/marcsantiago/gojay"
	"github.com/marcsantiago/gojay/benchmarks"
)

func BenchmarkJSONDecodeObjSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.SmallPayload{}
		_ = json.Unmarshal(benchmarks.SmallFixture, &result)
	}
}

func nothing(_ ...any) {}
func BenchmarkJSONParserSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = jsonparser.GetInt(benchmarks.SmallFixture, "tz")
		_, _ = jsonparser.GetInt(benchmarks.SmallFixture, "v")
		_, _ = jsonparser.GetInt(benchmarks.SmallFixture, "sid")
		_, _ = jsonparser.GetInt(benchmarks.SmallFixture, "st")
		_, _ = jsonparser.GetInt(benchmarks.SmallFixture, "gr")
		_, _, _, _ = jsonparser.Get(benchmarks.SmallFixture, "uuid")
		_, _, _, _ = jsonparser.Get(benchmarks.SmallFixture, "ua")
		nothing()
	}
}

func BenchmarkJsonIterDecodeObjSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.SmallPayload{}
		_ = jsoniter.Unmarshal(benchmarks.SmallFixture, &result)
	}
}

func BenchmarkEasyJsonDecodeObjSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.SmallPayload{}
		_ = easyjson.Unmarshal(benchmarks.SmallFixture, &result)
	}
}

func BenchmarkGoJayDecodeObjSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.SmallPayload{}
		_ = gojay.UnmarshalJSONObject(benchmarks.SmallFixture, &result)
	}
}

func BenchmarkGoJayUnsafeDecodeObjSmall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result := benchmarks.SmallPayload{}
		_ = gojay.Unsafe.UnmarshalJSONObject(benchmarks.SmallFixture, &result)
	}
}
