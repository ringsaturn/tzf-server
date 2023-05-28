package handler

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
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
		if name == key.Name || name == "all" {
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
	if name == "all" {
		return convert.Revert(d.data), nil
	}
	for _, shape := range d.data.Timezones {
		if shape.Name == name {
			return convert.RevertItem(shape), nil
		}
	}
	return nil, errors.New("not found")
}

var (
	finder tzf.F
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

func setupFuzzyFinder(dataPath string) (tzf.F, VisualizableTimezoneData, error) {
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

func setupPolygonFinder(dataPath string) (tzf.F, VisualizableTimezoneData, error) {
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

func Setup(option *SetupFinderOptions) *server.Hertz {
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

func setupEngine() *server.Hertz {
	h := server.Default(server.WithHostPorts(":8080"))

	h.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "template/*.html")))

	h.GET("/", Index)
	h.GET("/ping", Ping)

	apiV1Group := h.Group("/api/v1")
	apiV1Group.GET("/tz", GetTimezoneName)
	apiV1Group.GET("/tzs", GetTimezoneNames)
	apiV1Group.GET("/tzs/all", GetAllSupportTimezoneNames)
	apiV1Group.GET("/tz/geojson", GetTimezoneShape)

	webPageGroup := h.Group("/web")
	webPageGroup.GET("/tz", GetTimezoneInfoPage)
	webPageGroup.GET("/tzs", GetTimezonesInfoPage)
	webPageGroup.GET("/tzs/all", GetAllSupportTimezoneNamesPage)
	webPageGroup.GET("/tz/geojson/viewer", GetGeoJSONViewerForTimezone)

	return h
}

func Index(ctx context.Context, c *app.RequestContext) {
	c.Redirect(http.StatusTemporaryRedirect, []byte("/web/tzs/all"))
}

func Ping(ctx context.Context, c *app.RequestContext) { c.JSON(http.StatusOK, utils.H{}) }
