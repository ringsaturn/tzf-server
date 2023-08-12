package handler_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/redis/go-redis/v9"
	"github.com/ringsaturn/tzf-server/handler"
)

var (
	redisServerOnce = sync.Once{}

	redisClient = func() *redis.Client {
		redisOpt, err := redis.ParseURL("redis://localhost:6380")
		if err != nil {
			panic(err)
		}
		client := redis.NewClient(redisOpt)
		return client
	}()
)

func mustStartServer() {
	redisServerOnce.Do(func() {
		go func() { _ = handler.StartRedisServer() }()
		time.Sleep(100 * time.Millisecond)
	})
}

func TestRedisServerCallByClientGetTimezoneName(t *testing.T) {
	mustStartServer()

	ctx := context.Background()

	var testCases = []struct {
		lng float64
		lat float64
		tz  string
	}{
		{116.3883, 39.9289, "Asia/Shanghai"},
	}

	for _, case_ := range testCases {
		tz, err := redisClient.Do(ctx, "get_tz", case_.lng, case_.lat).Result()
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, case_.tz == tz)
	}
}

func TestRedisServerCallByClientGetTimezoneNames(t *testing.T) {
	mustStartServer()

	ctx := context.Background()

	var testCases = []struct {
		lng float64
		lat float64
		tz  []string
	}{
		{116.3883, 39.9289, []string{"Asia/Shanghai"}},
	}

	for _, case_ := range testCases {
		_tz, err := redisClient.Do(ctx, "get_tzs", case_.lng, case_.lat).Result()
		if err != nil {
			t.Fatal(err)
		}
		tz := _tz.([]interface{})
		assert.True(t, len(case_.tz) == len(tz))
		for i := 0; i < len(tz); i++ {
			assert.True(t, case_.tz[i] == tz[i])
		}
	}
}
