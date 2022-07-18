package main

import (
	"embed"
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

//go:embed static/*
var f embed.FS

type Request struct {
	Lng  float64 `form:"lng"`
	Lat  float64 `form:"lat"`
	Name string  `form:"name"`
}

func (req *Request) GetTimezoneData(finder *tzf.Finder) (*pb.Timezone, error) {
	if req.Name != "" {
		return finder.GetTimezoneShapeByName(req.Name)
	}
	return finder.GetTimezone(req.Lng, req.Lat)
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

	tz, tzErr := req.GetTimezoneData(h.finder)
	if tzErr != nil {
		ctx.String(500, tzErr.Error())
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
	tz, tzErr := req.GetTimezoneData(h.finder)
	if tzErr != nil {
		ctx.String(500, tzErr.Error())
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

	tz, tzErr := req.GetTimezoneData(h.finder)
	if tzErr != nil {
		ctx.String(500, tzErr.Error())
		return
	}

	dataAPIURL := fmt.Sprintf("http://%v/tz/geojson?lng=%v&lat=%v&name=%v", ctx.Request.Host, req.Lng, req.Lat, req.Name)
	geoJSONURL := fmt.Sprintf("http://geojson.io/#data=data:text/x-url,%v", url.QueryEscape(dataAPIURL))

	ctx.HTML(200, "info.html", gin.H{
		"Title":   tz.Name,
		"URL":     template.URL(geoJSONURL),
		"URLName": fmt.Sprintf("View Polygon for %v", tz.Name),
	})
}

type RequestByOffset struct {
	Offset int `form:"offset"`
}

type ResponseItem struct {
	Name string
	URL  string
}

type TZInfoPageResponse struct {
	Title string
	Items []*ResponseItem
}

func (h *Handler) GetTZInfoPageByOffset(ctx *gin.Context) {
	req := &RequestByOffset{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.String(400, err.Error())
		return
	}
	tzs, err := h.finder.GetTimezoneShapeByShift(req.Offset)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	resp := &TZInfoPageResponse{
		Title: fmt.Sprintf("tz for %v", req.Offset),
		Items: make([]*ResponseItem, 0),
	}

	for _, tz := range tzs {
		dataAPIURL := fmt.Sprintf("http://%v/tz/geojson?name=%v", ctx.Request.Host, tz.Name)
		geoJSONURL := fmt.Sprintf("http://geojson.io/#data=data:text/x-url,%v", url.QueryEscape(dataAPIURL))
		resp.Items = append(resp.Items, &ResponseItem{
			Name: tz.Name,
			URL:  geoJSONURL,
		})
	}

	ctx.HTML(200, "info_multi.html", resp)
}

func NewServer(finder *tzf.Finder) *gin.Engine {
	e := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "static/*.html"))
	e.SetHTMLTemplate(templ)
	h := Handler{
		finder: finder,
	}
	e.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	e.GET("/info", h.TZInfoPage)
	e.GET("/tz", h.GetTZ)
	e.GET("/tz/offset", h.GetTZInfoPageByOffset)
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
	panic(e.Run(*addr))
}
