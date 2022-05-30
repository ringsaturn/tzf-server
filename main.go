package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ringsaturn/tzf"
	tzfrel "github.com/ringsaturn/tzf-rel"
	"github.com/ringsaturn/tzf/convert"
	"github.com/ringsaturn/tzf/pb"
	"google.golang.org/protobuf/proto"
)

type TzRequest struct {
	Lng float64 `form:"lng"`
	Lat float64 `form:"lat"`
}

func main() {

	input := &pb.Timezones{}

	// Lite data, about 16.7MB
	dataFile := tzfrel.LiteData

	// Full data, about 83.5MB
	// dataFile := tzfrel.FullData

	if err := proto.Unmarshal(dataFile, input); err != nil {
		panic(err)
	}
	finder, _ := tzf.NewFinderFromPB(input)

	e := gin.Default()
	e.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	e.GET("/tz", func(ctx *gin.Context) {
		req := &TzRequest{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.String(400, err.Error())
			return
		}
		tz, err := finder.GetTimezone(req.Lng, req.Lat)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}
		ctx.String(200, tz.Name)
	})
	e.GET("/tz/geojson", func(ctx *gin.Context) {
		req := &TzRequest{}
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.String(400, err.Error())
			return
		}
		tz, err := finder.GetTimezone(req.Lng, req.Lat)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}
		ctx.JSON(200, convert.RevertItem(tz))
	})
	e.Run()
}
