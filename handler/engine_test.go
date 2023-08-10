package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/ringsaturn/tzf-server/handler"
)

var (
	h      = handler.Setup(nil)
	hFuzzy = handler.Setup(&handler.SetupFinderOptions{FinderType: handler.FuzzyFinder})
)

func TestIndex(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/ping", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestPing(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/ping", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestGetTimezoneName(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tz?lng=116.3883&lat=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &handler.GetTimezoneNameResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.True(t, result.Timezone == "Asia/Shanghai")
}

func TestFuzzyGetTimezoneName(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tz?lng=116.3883&lat=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &handler.GetTimezoneNameResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.True(t, result.Timezone == "Asia/Shanghai")
}

func TestFuzzyGetTimezoneNames(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tzs?lng=87.6168&lat=43.8254", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &handler.GetTimezoneNamesResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.DeepEqual(t, result.Timezones, []string{"Asia/Shanghai", "Asia/Urumqi"})
}

// func TestWebGetAllTimezoneNames(t *testing.T) {
// 	w := ut.PerformRequest(h.Engine, "GET", "/web/tzs/all", nil)
// 	resp := w.Result()
// 	fmt.Println(string(resp.BodyBytes()))
// 	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
// }

// func TestWebGetTimezoneName(t *testing.T) {
// 	w := ut.PerformRequest(h.Engine, "GET", "/web/tz?lng=116.3883&lat=39.9289", nil)
// 	resp := w.Result()
// 	fmt.Println(string(resp.BodyBytes()))
// 	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
// }
