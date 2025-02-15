package httpserver_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/google/go-cmp/cmp"
	v1 "github.com/ringsaturn/tzf-server/gen/go/tzf_server/v1"
	"github.com/ringsaturn/tzf-server/internal/config"
	"github.com/ringsaturn/tzf-server/internal/finder"
	"github.com/ringsaturn/tzf-server/internal/httpserver"
	"github.com/ringsaturn/tzf-server/internal/wraps"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
)

func setUp(cfg *config.Config) *server.Hertz {
	tZfinder, err := finder.NewFinder(cfg)
	if err != nil {
		panic(err)
	}
	tzfServiceHTTPServer := wraps.NewTZFServiceServer(tZfinder)
	webHandler := httpserver.NewWebHandler(tZfinder)
	hertz := httpserver.NewServer(cfg, tzfServiceHTTPServer, webHandler)
	return hertz
}

var (
	h *server.Hertz = setUp(&config.Config{
		HertzPrometheusHostPorts: "localhost:7777",
		HertzPrometheusPath:      "/metrics",
	})
	hFuzzy = setUp(&config.Config{
		FinderType:               int(finder.FuzzyFinder),
		HertzPrometheusHostPorts: "localhost:7778",
		HertzPrometheusPath:      "/metrics2",
	})
)

func mustEqualForProto(t *testing.T, expected proto.Message, actual proto.Message) {
	eq := proto.Equal(expected, actual)
	if !eq {
		diff := cmp.Diff(expected, actual, protocmp.Transform())
		t.Fatal(diff)
	}
}

func TestRoot(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusTemporaryRedirect, resp.StatusCode())
}

func TestPing(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/ping", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestGetTimezoneName(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tz?longitude=116.3883&latitude=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &v1.GetTimezoneResponse{}
	err := protojson.Unmarshal(resp.BodyBytes(), result)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &v1.GetTimezoneResponse{
		Timezone:     "Asia/Shanghai",
		Abbreviation: "CST",
		Offset:       durationpb.New(28800 * time.Second),
	}
	mustEqualForProto(t, expected, result)
}

func TestFuzzyGetTimezoneName(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tz?longitude=116.3883&latitude=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &v1.GetTimezoneResponse{}
	err := protojson.Unmarshal(resp.BodyBytes(), result)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &v1.GetTimezoneResponse{
		Timezone:     "Asia/Shanghai",
		Abbreviation: "CST",
		Offset:       durationpb.New(28800 * time.Second),
	}
	mustEqualForProto(t, expected, result)
}

func TestFuzzyGetTimezoneNames(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tzs?longitude=87.6168&latitude=43.8254", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())

	result := &v1.GetTimezonesResponse{}
	err := protojson.Unmarshal(resp.BodyBytes(), result)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := &v1.GetTimezonesResponse{
		Timezones: []*v1.GetTimezoneResponse{
			{
				Timezone:     "Asia/Shanghai",
				Abbreviation: "CST",
				Offset:       durationpb.New(28800 * time.Second),
			},
			{
				Timezone:     "Asia/Urumqi",
				Abbreviation: "+06",
				Offset:       durationpb.New(21600 * time.Second),
			},
		},
	}
	mustEqualForProto(t, expected, result)
}

func TestGetTimezoneShape(t *testing.T) {
	w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tz/geojson?longitude=116.3883&latitude=39.9289", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestFuzzyGetTimezoneShape(t *testing.T) {
	w := ut.PerformRequest(hFuzzy.Engine, "GET", "/api/v1/tz/geojson?longitude=116.3883&latitude=39.9289", nil)
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
// 	w := ut.PerformRequest(h.Engine, "GET", "/web/tz?longitude=116.3883&latitude=39.9289", nil)
// 	resp := w.Result()
// 	fmt.Println(string(resp.BodyBytes()))
// 	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
// }
