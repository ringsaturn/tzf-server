package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type LocationRequest struct {
	Lng float64 `query:"lng"`
	Lat float64 `query:"lat"`
}

type GetTimezoneNameResponse struct {
	Timezone string `json:"timezone"`
}

func GetTimezoneName(ctx context.Context, c *app.RequestContext) {
	req := &LocationRequest{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}
	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}
	c.JSON(http.StatusOK, &GetTimezoneNameResponse{Timezone: timezone})
}

type GetTimezoneNamesResponse struct {
	Timezones []string `json:"timezones"`
}

func GetTimezoneNames(ctx context.Context, c *app.RequestContext) {
	req := &LocationRequest{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}
	timezones, err := finder.GetTimezoneNames(req.Lng, req.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &GetTimezoneNamesResponse{Timezones: timezones})
}

func GetAllSupportTimezoneNames(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{"timezones": finder.TimezoneNames()})
}

type GetTimezoneInfoRequest struct {
	Name string  `query:"name"`
	Lng  float64 `query:"lng"`
	Lat  float64 `query:"lat"`
}

func GetTimezoneShape(ctx context.Context, c *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	err := c.Bind(req)
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
