package handler

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb/maptile"
	"github.com/ringsaturn/tzf"
	tzfrel "github.com/ringsaturn/tzf-rel"
	"github.com/ringsaturn/tzf/convert"
	"github.com/ringsaturn/tzf/pb"
	"github.com/ringsaturn/tzf/reduce"
	"google.golang.org/protobuf/proto"
)

//go:embed template/*
var f embed.FS

//go:embed static/*
var staticFS embed.FS

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Finder interface {
	GetTimezoneName(lng float64, lat float64) string
	GetTimezoneNames(lng float64, lat float64) ([]string, error)
	TimezoneNames() []string
}

type VisualizableTimezoneData interface {
	GetShape(name string) (interface{}, error)
}

type fuzzyData struct {
	data *pb.PreindexTimezones
}

func (d *fuzzyData) GetShape(name string) (interface{}, error) {
	var hit bool
	tileSet := maptile.Set{}
	for _, key := range d.data.Keys {
		if name == key.Name {
			hit = true
			tileSet[maptile.New(uint32(key.X), uint32(key.Y), maptile.Zoom(key.Z))] = true
		}
	}
	if !hit {
		return nil, errors.New("not found")
	}
	return tileSet.ToFeatureCollection(), nil
}

type polygonData struct {
	data *pb.Timezones
}

func (d *polygonData) GetShape(name string) (interface{}, error) {
	for _, shape := range d.data.Timezones {
		if shape.Name == name {
			return convert.RevertItem(shape), nil
		}
	}
	return nil, errors.New("not found")
}

//

var (
	finder Finder
	tzData VisualizableTimezoneData
)

type FinderType int

const (
	PolygonFinder FinderType = iota
	FuzzyFinder
)

type SetupFinderOptions struct {
	FinderType     FinderType
	CustomDataPath string
}

//

func setupFuzzyFinder(dataPath string) (Finder, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.PreindexTimezones{}
	if dataPath == "" {
		err = proto.Unmarshal(tzfrel.PreindexData, tzpb)
		if err != nil {
			return nil, nil, err
		}
	} else {
		rawFile, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, nil, err
		}
		err = proto.Unmarshal(rawFile, tzpb)
		if err != nil {
			return nil, nil, err
		}
	}
	finder, err = tzf.NewFuzzyFinderFromPB(tzpb)
	return finder, &fuzzyData{data: tzpb}, err
}

func setupPolygonFinder(dataPath string) (Finder, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.Timezones{}

	if dataPath == "" {
		compressPb := &pb.CompressedTimezones{}
		err = proto.Unmarshal(tzfrel.LiteCompressData, compressPb)
		if err != nil {
			return nil, nil, err
		}
		tzpb, err = reduce.Decompress(compressPb)
		if err != nil {
			return nil, nil, err
		}
	} else {
		rawFile, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, nil, err
		}
		err = proto.Unmarshal(rawFile, tzpb)
		if err != nil {
			return nil, nil, err
		}
	}
	finder, err = tzf.NewFinderFromPB(tzpb)
	return finder, &polygonData{data: tzpb}, err
}

func Setup(option *SetupFinderOptions) *gin.Engine {
	if option == nil {
		option = &SetupFinderOptions{
			FinderType: PolygonFinder,
		}
	}
	var err error
	switch option.FinderType {
	case FuzzyFinder:
		finder, tzData, err = setupFuzzyFinder(option.CustomDataPath)
	default:
		finder, tzData, err = setupPolygonFinder(option.CustomDataPath)
	}
	check(err)
	return setupEngine()
}

func setupEngine() *gin.Engine {
	engine := gin.Default()

	templates := template.Must(template.New("").ParseFS(f, "template/*.html"))
	engine.SetHTMLTemplate(templates)

	fe, _ := fs.Sub(staticFS, "static")
	engine.StaticFS("/static", http.FS(fe))

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://geojson.io", "http://geojson.io"},
		AllowMethods:     []string{"PUT", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/ping", Ping)
	engine.GET("/tz", GetTimezoneName)
	engine.GET("/tz/info", GetTimezoneInfoPage)
	engine.GET("/tzs", GetTimezoneNames)
	engine.GET("/tzs/list", GetAllSupportTimezoneNames)
	engine.GET("/tzs/list/page", GetAllSupportTimezoneNamesPage)
	engine.GET("/tz/geojson", GetTimezoneShape)
	engine.GET("/tz/geojson/viewer", GetGeoJSONViewerForTimezone)
	return engine
}

func Ping(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }

type LocationRequest struct {
	Lng float64 `form:"lng"`
	Lat float64 `form:"lat"`
}

func GetTimezoneName(c *gin.Context) {
	req := &LocationRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}
	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no timezone found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"timezone": timezone})
}

func GetTimezoneNames(c *gin.Context) {
	req := &LocationRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}
	timezones, err := finder.GetTimezoneNames(req.Lng, req.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"timezones": timezones})
}

func GetAllSupportTimezoneNames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"timezones": finder.TimezoneNames()})
}

type GetTimezoneInfoRequest struct {
	Name string  `form:"name"`
	Lng  float64 `form:"lng"`
	Lat  float64 `form:"lat"`
}

func GetTimezoneShape(c *gin.Context) {
	req := &GetTimezoneInfoRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = finder.GetTimezoneName(req.Lng, req.Lat)
	}
	shape, err := tzData.GetShape(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shape)
}

func GetTimezoneInfoPage(c *gin.Context) {
	req := &GetTimezoneInfoRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}

	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no timezone found"})
		return
	}

	viewerURL := fmt.Sprintf("http://%v/tz/geojson/viewer?lng=%v&lat=%v&name=%v", string(c.Request.Host), req.Lng, req.Lat, req.Name)

	c.HTML(200, "info.html", map[string]any{
		"Title": timezone,
		"URL":   template.URL(viewerURL),
		"Name":  fmt.Sprintf("View Polygon for %v", timezone),
	})
}

type GetTimezoneInfoPageResponseItem struct {
	Name string
	URL  string
}

type GetTimezoneInfoPageResponse struct {
	Title string
	Items []*GetTimezoneInfoPageResponseItem
}

func GetAllSupportTimezoneNamesPage(c *gin.Context) {
	resp := &GetTimezoneInfoPageResponse{
		Title: "All timezones",
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	for _, name := range finder.TimezoneNames() {
		viewerURL := fmt.Sprintf("http://%v/tz/geojson/viewer?name=%v", string(c.Request.Host), name)
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  viewerURL,
		})
	}
	c.HTML(200, "info_multi.html", resp)
}

func GetGeoJSONViewerForTimezone(c *gin.Context) {
	req := &GetTimezoneInfoRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}

	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no timezone found"})
		return
	}

	url := fmt.Sprintf("http://%v/tz/geojson?lng=%v&lat=%v&name=%v", c.Request.Host, req.Lng, req.Lat, req.Name)
	c.HTML(http.StatusOK, "viewer.html", map[string]any{
		"URL": url,
	})
}
