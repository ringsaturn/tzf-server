package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ringsaturn/tzf"
	tzfrel "github.com/ringsaturn/tzf-rel"
	"github.com/ringsaturn/tzf/convert"
	"github.com/ringsaturn/tzf/pb"
	"google.golang.org/protobuf/proto"
)

type Request struct {
	Lng float64 `form:"lng"`
	Lat float64 `form:"lat"`
}

type Handler struct {
	finder *tzf.Finder
}

func (h *Handler) GetTZ(ctx *gin.Context) {
	req := &Request{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.String(400, err.Error())
		return
	}
	tz, err := h.finder.GetTimezone(req.Lng, req.Lat)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}
	ctx.String(200, tz.Name)
}

func (h *Handler) GetTZGeoJSON(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	req := &Request{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.String(400, err.Error())
		return
	}
	tz, err := h.finder.GetTimezone(req.Lng, req.Lat)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}
	ctx.JSON(200, convert.RevertItem(tz))
}

func (h *Handler) TZInfoPage(ctx *gin.Context) {
	req := &Request{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.String(400, err.Error())
		return
	}
	tz, err := h.finder.GetTimezone(req.Lng, req.Lat)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	dataAPIURL := fmt.Sprintf("http://%v/tz/geojson?lng=%v&lat=%v", ctx.Request.Host, req.Lng, req.Lat)
	geoJSONURL := fmt.Sprintf("http://geojson.io/#data=data:text/x-url,%v", url.QueryEscape(dataAPIURL))

	ctx.HTML(200, "info.html", gin.H{
		"Title":   tz.Name,
		"URL":     template.URL(geoJSONURL),
		"URLName": fmt.Sprintf("View Polygon for %v", tz.Name),
	})
}

func NewServer(finder *tzf.Finder) *gin.Engine {
	e := gin.Default()
	e.LoadHTMLFiles("info.html")
	h := Handler{
		finder: finder,
	}
	e.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	e.GET("/info", h.TZInfoPage)
	e.GET("/tz", h.GetTZ)
	e.GET("/tz/geojson", h.GetTZGeoJSON)
	return e
}

func NewTZData(tzpbPath string) []byte {
	if tzpbPath == "" {
		return tzfrel.LiteData
	}

	rawFile, err := ioutil.ReadFile(tzpbPath)
	if err != nil {
		panic(err)
	}
	return rawFile

}

func main() {
	addr := flag.String("addr", ":8080", "API Server Addr")
	tzpbPath := flag.String("tzpb", "", "custom tzpb data path. Use `tzfrel.LiteData` by default.")
	flag.Parse()

	input := &pb.Timezones{}
	if err := proto.Unmarshal(NewTZData(*tzpbPath), input); err != nil {
		panic(err)
	}
	finder, _ := tzf.NewFinderFromPB(input)

	e := NewServer(finder)
	e.Run(*addr)
}
