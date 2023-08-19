package handler_test

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
)

func BenchmarkGetTimezoneName_Beijing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ut.PerformRequest(h.Engine, "GET", "/api/v1/tz?lng=116.3883&lat=39.9289", nil)
	}
}

func BenchmarkGetTimezoneName_Urumqi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ut.PerformRequest(h.Engine, "GET", "/api/v1/tzs?lng=87.6168&lat=43.8254", nil)
	}
}
