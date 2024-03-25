package wraps

import (
	"context"
	"errors"

	"github.com/google/wire"
	"github.com/ringsaturn/tzf-server/internal/finder"
	v1 "github.com/ringsaturn/tzf-server/tzf/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type apiServer struct {
	tzfinder *finder.TZfinder

	v1.TZFServiceServer
	v1.TZFServiceHTTPServer
}

func (as *apiServer) GetAllTimezones(ctx context.Context, in *v1.GetAllTimezonesRequest) (*v1.GetAllTimezonesResponse, error) {
	names := as.tzfinder.TimezoneNames()
	items := make([]*v1.GetTimezoneResponse, len(names))
	for i, timezone := range names {
		items[i] = &v1.GetTimezoneResponse{
			Timezone:     timezone,
			Abbreviation: as.tzfinder.TZName2Abbreviation[timezone],
			Offset:       as.tzfinder.TZName2Offset[timezone],
		}
	}
	return &v1.GetAllTimezonesResponse{Timezones: items}, nil
}

func (as *apiServer) GetTimezone(ctx context.Context, in *v1.GetTimezoneRequest) (*v1.GetTimezoneResponse, error) {
	zone := as.tzfinder.GetTimezoneName(in.Longitude, in.Latitude)
	abbr, ok := as.tzfinder.TZName2Abbreviation[zone]
	if !ok {
		return nil, errors.New("not found")
	}
	offset, ok := as.tzfinder.TZName2Offset[zone]
	if !ok {
		return nil, errors.New("not found")
	}
	return &v1.GetTimezoneResponse{
		Timezone:     zone,
		Abbreviation: abbr,
		Offset:       offset,
	}, nil
}

func (as *apiServer) GetTimezones(ctx context.Context, in *v1.GetTimezonesRequest) (*v1.GetTimezonesResponse, error) {
	zones, err := as.tzfinder.GetTimezoneNames(in.Longitude, in.Latitude)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.GetTimezoneResponse, len(zones))
	for i, timezone := range zones {
		items[i] = &v1.GetTimezoneResponse{
			Timezone:     timezone,
			Abbreviation: as.tzfinder.TZName2Abbreviation[timezone],
			Offset:       as.tzfinder.TZName2Offset[timezone],
		}
	}
	return &v1.GetTimezonesResponse{Timezones: items}, nil
}

func (as *apiServer) Ping(ctx context.Context, in *v1.PingRequest) (*v1.PingResponse, error) {
	return &v1.PingResponse{
		Data: &anypb.Any{
			TypeUrl: "tzf.v1.PingResponse",
			Value:   []byte("pong"),
		},
	}, nil
}

func NewTZFServiceServer(tzfinder *finder.TZfinder) v1.TZFServiceHTTPServer {
	return &apiServer{
		tzfinder: tzfinder,
	}
}

var HTTPProviderSet = wire.NewSet(
	NewTZFServiceServer,
)

// var GRPCProviderSet = wire.NewSet(
// 	NewTZFServiceServer,
// 	wire.Bind(new(v1.TZFServiceServer), new(*apiServer)),
// )
