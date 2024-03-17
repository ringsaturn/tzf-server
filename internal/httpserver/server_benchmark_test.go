package httpserver_test

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
)

func BenchmarkGetTimezoneName_Beijing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ut.PerformRequest(h.Engine, "GET", "/api/v1/tz?longitude=116.3883&latitude=39.9289", nil)
	}
}

func BenchmarkGetTimezoneName_Urumqi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ut.PerformRequest(h.Engine, "GET", "/api/v1/tzs?longitude=87.6168&latitude=43.8254", nil)
	}
}
