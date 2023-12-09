package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/google/go-cmp/cmp"
	"github.com/ringsaturn/tzf-server/internal/handler"
	v1 "github.com/ringsaturn/tzf-server/proto/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	h      = handler.Setup(zap.Must(zap.NewProduction()), nil)
	hFuzzy = handler.Setup(zap.Must(zap.NewProduction()), &handler.SetupFinderOptions{FinderType: handler.FuzzyFinder})
)

func mustEqualForProto(t *testing.T, expected proto.Message, actual proto.Message) {
	eq := proto.Equal(expected, actual)
	if !eq {
		diff := cmp.Diff(expected, actual)
		t.Fatalf(diff)
	}
}

func TestRoot(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusTemporaryRedirect, resp.StatusCode())
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

	result := &v1.GetTimezoneResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := &v1.GetTimezoneResponse{
		Timezone:     "Asia/Shanghai",
		Abbreviation: "CST",
		Offset:       28800,
	}
	mustEqualForProto(t, expected, result)
}

func TestFuzzyGetTimezoneName(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tz?lng=116.3883&lat=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &v1.GetTimezoneResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := &v1.GetTimezoneResponse{
		Timezone:     "Asia/Shanghai",
		Abbreviation: "CST",
		Offset:       28800,
	}
	mustEqualForProto(t, expected, result)
}

func TestFuzzyGetTimezoneNames(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tzs?lng=87.6168&lat=43.8254", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &v1.GetTimezonesResponse{}
	err := json.Unmarshal(resp.BodyBytes(), &result)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := &v1.GetTimezonesResponse{
		Timezones: []*v1.GetTimezoneResponse{
			{
				Timezone:     "Asia/Shanghai",
				Abbreviation: "CST",
				Offset:       28800,
			},
			{
				Timezone:     "Asia/Urumqi",
				Abbreviation: "+06",
				Offset:       21600,
			},
		},
	}
	mustEqualForProto(t, expected, result)
}

func TestGetTimezoneShape(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tz/geojson?lng=116.3883&lat=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestFuzzyGetTimezoneShape(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tz/geojson?lng=116.3883&lat=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestGetAllSupportTimezoneNames(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tzs/all", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
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
