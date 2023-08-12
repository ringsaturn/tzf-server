package handler_test

import (
	"testing"

	"github.com/ringsaturn/tzf-server/handler"
	"github.com/tidwall/redcon"
	gomock "go.uber.org/mock/gomock"
)

func TestRedisServerGetTimezoneName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	conn.EXPECT().WriteString("Asia/Shanghai").MaxTimes(1)

	cmd := redcon.Command{
		Raw: []byte("get_tz 116.3883 39.9289"),
		Args: [][]byte{
			[]byte("get_tz"),
			[]byte("116.3883"),
			[]byte("39.9289"),
		},
	}
	handler.RedisGetTZCmd(conn, cmd)
}

func BenchmarkRedisServerGetTimezoneName(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	conn.EXPECT().WriteString("Asia/Shanghai").MinTimes(1)

	cmd := redcon.Command{
		Raw: []byte("get_tz 116.3883 39.9289"),
		Args: [][]byte{
			[]byte("get_tz"),
			[]byte("116.3883"),
			[]byte("39.9289"),
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.RedisGetTZCmd(conn, cmd)
	}
}

func TestRedisServerGetTimezoneNameWithInvalidArgs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	conn.EXPECT().WriteError(gomock.Any()).MaxTimes(1)

	cmd := redcon.Command{
		Raw: []byte("get_tz 116.3883"),
		Args: [][]byte{
			[]byte("get_tz"),
			[]byte("116.3883"),
		},
	}
	handler.RedisGetTZCmd(conn, cmd)
}

func TestRedisServerGetTimezoneNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	conn.EXPECT().WriteArray(2).MaxTimes(1)
	conn.EXPECT().WriteBulkString("Asia/Shanghai").MaxTimes(1)
	conn.EXPECT().WriteBulkString("Asia/Urumqi").MaxTimes(1)

	cmd := redcon.Command{
		Raw: []byte("get_tzs 87.6168 43.8254"),
		Args: [][]byte{
			[]byte("get_tzs"),
			[]byte("87.6168"),
			[]byte("43.8254"),
		},
	}
	handler.RedisGetTZsCmd(conn, cmd)
}

func BenchmarkRedisServerGetTimezoneNames(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	conn.EXPECT().WriteArray(2).MinTimes(1)
	conn.EXPECT().WriteBulkString("Asia/Shanghai").MinTimes(1)
	conn.EXPECT().WriteBulkString("Asia/Urumqi").MinTimes(1)

	cmd := redcon.Command{
		Raw: []byte("get_tzs 87.6168 43.8254"),
		Args: [][]byte{
			[]byte("get_tzs"),
			[]byte("87.6168"),
			[]byte("43.8254"),
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.RedisGetTZsCmd(conn, cmd)
	}
}
