package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	v1 "github.com/ringsaturn/tzf-server/proto/v1"
)

func GetTimezoneName(ctx context.Context, in *v1.GetTimezoneRequest) (*v1.GetTimezoneResponse, error) {
	timezone := finder.GetTimezoneName(in.Longitude, in.Latitude)
	if timezone == "" {
		return nil, errors.New("no timezone found")
	}
	return &v1.GetTimezoneResponse{
		Timezone:     timezone,
		Abbreviation: tzName2Abbreviation[timezone],
		Offset:       tzName2Offset[timezone],
	}, nil
}

func GetTimezoneNames(ctx context.Context, in *v1.GetTimezonesRequest) (*v1.GetTimezonesResponse, error) {
	timezones, err := finder.GetTimezoneNames(in.Longitude, in.Latitude)
	if err != nil {
		return nil, fmt.Errorf("failed to get timezone names: %w", err)
	}
	items := make([]*v1.GetTimezoneResponse, len(timezones))
	for i, timezone := range timezones {
		items[i] = &v1.GetTimezoneResponse{
			Timezone:     timezone,
			Abbreviation: tzName2Abbreviation[timezone],
			Offset:       tzName2Offset[timezone],
		}
	}
	return &v1.GetTimezonesResponse{Timezones: items}, nil
}

func GetAllSupportTimezoneNames(ctx context.Context, in *v1.GetAllTimezonesRequest) (*v1.GetAllTimezonesResponse, error) {
	items := make([]*v1.GetTimezoneResponse, len(finder.TimezoneNames()))
	for i, timezone := range finder.TimezoneNames() {
		items[i] = &v1.GetTimezoneResponse{
			Timezone:     timezone,
			Abbreviation: tzName2Abbreviation[timezone],
			Offset:       tzName2Offset[timezone],
		}
	}
	return &v1.GetAllTimezonesResponse{Timezones: items}, nil
}

type GetTimezoneInfoRequest struct {
	Name string  `query:"name"`
	Lng  float64 `query:"lng"`
	Lat  float64 `query:"lat"`
}

func GetTimezoneShape(ctx context.Context, c *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"err": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = finder.GetTimezoneName(req.Lng, req.Lat)
	}
	shape, err := tzData.GetShape(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shape)
}
