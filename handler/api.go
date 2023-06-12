package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/pkg/errors"
	"github.com/tidwall/redcon"
)

type LocationRequest struct {
	Lng float64 `query:"lng"`
	Lat float64 `query:"lat"`
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
	c.JSON(http.StatusOK, utils.H{"timezone": timezone})
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
	c.JSON(http.StatusOK, utils.H{"timezones": timezones})
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

func parseCoordinates(cmd redcon.Command) (float64, float64, error) {
	if len(cmd.Args) != 3 {
		return 0, 0, errors.New("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
	}
	lng, err := strconv.ParseFloat(string(cmd.Args[1]), 64)
	if err != nil {
		return 0, 0, err
	}
	lat, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		return 0, 0, err
	}
	return lng, lat, nil
}

func RedisGetTZCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_name := finder.GetTimezoneName(lng, lat)
	if timezone_name == "" {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteString(timezone_name)
}

func RedisGetTZsCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_names, err := finder.GetTimezoneNames(lng, lat)
	if err != nil {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteArray(len(timezone_names))
	for _, name := range timezone_names {
		conn.WriteBulkString(name)
	}
}
