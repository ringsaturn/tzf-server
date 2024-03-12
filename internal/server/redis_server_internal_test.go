package server

import (
	"testing"

	"github.com/tidwall/redcon"
)

func Test_parseCoordinates(t *testing.T) {
	type args struct {
		cmd redcon.Command
	}
	tests := []struct {
		name    string
		args    args
		lng     float64
		lat     float64
		wantErr bool
	}{
		{
			"test1",
			args{
				redcon.Command{
					Args: [][]byte{
						[]byte("get_tz"),
						[]byte("116.3883"),
						[]byte("39.9289"),
					}}},
			116.3883, 39.9289, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseCoordinates(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.lng {
				t.Errorf("parseCoordinates() got = %v, want %v", got, tt.lng)
			}
			if got1 != tt.lat {
				t.Errorf("parseCoordinates() got1 = %v, want %v", got1, tt.lat)
			}
		})
	}
}

func Benchmark_parseCoordinates(b *testing.B) {
	cmd := redcon.Command{
		Args: [][]byte{
			[]byte("get_tz"),
			[]byte("116.3883"),
			[]byte("39.9289"),
		}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parseCoordinates(cmd)
	}
}
