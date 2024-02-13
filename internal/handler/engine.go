package handler

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/paulmach/orb/maptile"
	"github.com/ringsaturn/tzf"
	tzfrellite "github.com/ringsaturn/tzf-rel-lite"
	v1 "github.com/ringsaturn/tzf-server/proto/v1"
	"github.com/ringsaturn/tzf/convert"
	"github.com/ringsaturn/tzf/pb"
	"github.com/ringsaturn/tzf/reduce"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
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
	finder              tzf.F
	tzData              VisualizableTimezoneData
	tzName2Abbreviation = make(map[string]string)
	tzName2Offset       = make(map[string]*durationpb.Duration)
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

func postSetUp(finder tzf.F) error {
	for _, tzName := range finder.TimezoneNames() {
		location, err := time.LoadLocation(tzName)

		if err != nil {
			return err
		}
		abbreviation, offset := time.Now().In(location).Zone()
		tzName2Abbreviation[tzName] = abbreviation
		tzName2Offset[tzName] = durationpb.New(time.Duration(offset * int(time.Second)))
	}
	return nil
}

func setupFuzzyFinder(logger *zap.Logger, dataPath string) (tzf.F, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.PreindexTimezones{}
	if dataPath == "" {
		logger.Debug("Fuzzy finder will use data from tzf-rel")
		err = proto.Unmarshal(tzfrellite.PreindexData, tzpb)
		if err != nil {
			logger.Sugar().Error("Unmarshal failed", "err", err)
			return nil, nil, err
		}
	} else {
		logger.Debug("Fuzzy finder use custom data")
		rawFile, err := os.ReadFile(dataPath)
		if err != nil {
			logger.Sugar().Error("Unable to load custom data", "err", err)
			return nil, nil, err
		}
		err = proto.Unmarshal(rawFile, tzpb)
		if err != nil {
			logger.Sugar().Error("Unmarshal failed", "err", err)
			return nil, nil, err
		}
	}
	finder, err = tzf.NewFuzzyFinderFromPB(tzpb)
	logger.Sugar().Debug("FuzzyFinder setup finished", "err", err)
	return finder, &fuzzyData{data: tzpb}, err
}

func setupPolygonFinder(logger *zap.Logger, dataPath string) (tzf.F, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.Timezones{}

	if dataPath == "" {
		logger.Debug("Finder will use data from tzf-rel")
		compressPb := &pb.CompressedTimezones{}
		err = proto.Unmarshal(tzfrellite.LiteCompressData, compressPb)
		if err != nil {
			return nil, nil, err
		}
		tzpb, err = reduce.Decompress(compressPb)
		if err != nil {
			return nil, nil, err
		}
	} else {
		logger.Debug("Finder will use data from tzf-rel")
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

func Setup(logger *zap.Logger, finderOption *SetupFinderOptions, cfg ...config.Option) *server.Hertz {
	if finderOption == nil {
		logger.Debug("option is nil, use default polygon finder")
		finderOption = &SetupFinderOptions{
			FinderType: PolygonFinder,
		}
	}
	var err error
	switch finderOption.FinderType {
	case FuzzyFinder:
		finder, tzData, err = setupFuzzyFinder(logger, finderOption.CustomDataPath)
	default:
		finder, tzData, err = setupPolygonFinder(logger, finderOption.CustomDataPath)
	}
	check(err)
	check(postSetUp(finder))
	hlog.SetLogger(hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WithFatalHook(zapcore.WriteThenPanic))))
	return setupEngine(cfg...)
}

type apiService struct{}

func (apiSrv *apiService) GetAllTimezones(ctx context.Context, in *v1.GetAllTimezonesRequest) (*v1.GetAllTimezonesResponse, error) {
	return GetAllSupportTimezoneNames(ctx, in)
}

func (apiSrv *apiService) GetTimezone(ctx context.Context, in *v1.GetTimezoneRequest) (*v1.GetTimezoneResponse, error) {
	return GetTimezoneName(ctx, in)
}

func (apiSrv *apiService) GetTimezones(ctx context.Context, in *v1.GetTimezonesRequest) (*v1.GetTimezonesResponse, error) {
	return GetTimezoneNames(ctx, in)
}

func setupEngine(cfg ...config.Option) *server.Hertz {
	h := server.New(cfg...)
	h.Use(
		recovery.Recovery(recovery.WithRecoveryHandler(
			func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
				c.JSON(http.StatusInternalServerError, utils.H{
					"error": fmt.Sprintf("[Recovery] err=%v\nstack=%s", err, stack),
				})
			},
		)),
	)

	h.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "template/*.html")))

	h.GET("/", Index)
	h.GET("/ping", Ping)

	v1.RegisterTZFServiceHTTPServer(h, &apiService{})

	h.GET("/api/v1/tz/geojson", GetTimezoneShape)

	webPageGroup := h.Group("/web")
	webPageGroup.GET("/tz", GetTimezoneInfoPage)
	webPageGroup.GET("/tzs", GetTimezonesInfoPage)
	webPageGroup.GET("/tzs/all", GetAllSupportTimezoneNamesPage)
	webPageGroup.GET("/tz/geojson/viewer", GetGeoJSONViewerForTimezone)
	webPageGroup.GET("/click", GetClickPage)

	return h
}

func Index(ctx context.Context, c *app.RequestContext) {
	c.Redirect(http.StatusTemporaryRedirect, []byte("/web/tzs/all"))
}

func Ping(ctx context.Context, c *app.RequestContext) { c.JSON(http.StatusOK, utils.H{}) }
