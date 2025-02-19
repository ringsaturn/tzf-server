package finder

import (
	"errors"
	"os"
	"time"

	"github.com/google/wire"
	"github.com/paulmach/orb/maptile"
	"github.com/ringsaturn/tzf"
	tzfrellite "github.com/ringsaturn/tzf-rel-lite"
	"github.com/ringsaturn/tzf-server/internal"
	"github.com/ringsaturn/tzf-server/internal/config"
	"github.com/ringsaturn/tzf/convert"
	pb "github.com/ringsaturn/tzf/gen/go/tzf/v1"
	"github.com/ringsaturn/tzf/reduce"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
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

type VisualizableTimezoneData interface {
	GetShape(name string) (interface{}, error)
}

type fuzzyData struct {
	data *pb.PreindexTimezones
}

func (d *fuzzyData) GetShape(name string) (interface{}, error) {
	feature := &convert.FeatureItem{
		Type: convert.FeatureType,
		Properties: convert.PropertiesDefine{
			Tzid: name,
		},
		Geometry: convert.GeometryDefine{
			Type:        convert.MultiPolygonType,
			Coordinates: make(convert.MultiPolygonCoordinates, 0),
		},
	}

	res := convert.MultiPolygonCoordinates{}

	var hit bool
	for _, key := range d.data.Keys {
		if name == key.Name || name == "all" {
			hit = true

			tile := maptile.New(uint32(key.X), uint32(key.Y), maptile.Zoom(key.Z))
			bound := tile.Bound()

			res = append(
				res,
				convert.PolygonCoordinates{
					{
						{bound.Min[0], bound.Min[1]},
						{bound.Max[0], bound.Min[1]},
						{bound.Max[0], bound.Max[1]},
						{bound.Min[0], bound.Max[1]},
					},
				})
		}
	}
	if !hit {
		return nil, errors.New("not found")
	}
	feature.Geometry.Coordinates = res

	return feature, nil
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

type F interface {
	tzf.F
}

type TZfinder struct {
	Finder              tzf.F
	TZData              VisualizableTimezoneData
	TZName2Abbreviation map[string]string
	TZName2Offset       map[string]*durationpb.Duration
}

func setupFuzzyFinder(dataPath string) (tzf.F, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.PreindexTimezones{}
	if dataPath == "" {
		internal.Loggger.Debug("Fuzzy finder will use data from tzf-rel")
		err = proto.Unmarshal(tzfrellite.PreindexData, tzpb)
		if err != nil {
			internal.Loggger.Sugar().Error("Unmarshal failed", "err", err)
			return nil, nil, err
		}
	} else {
		internal.Loggger.Debug("Fuzzy finder use custom data")
		rawFile, err := os.ReadFile(dataPath)
		if err != nil {
			internal.Loggger.Sugar().Error("Unable to load custom data", "err", err)
			return nil, nil, err
		}
		err = proto.Unmarshal(rawFile, tzpb)
		if err != nil {
			internal.Loggger.Sugar().Error("Unmarshal failed", "err", err)
			return nil, nil, err
		}
	}
	finder, err := tzf.NewFuzzyFinderFromPB(tzpb)
	internal.Loggger.Sugar().Debug("FuzzyFinder setup finished", "err", err)
	return finder, &fuzzyData{data: tzpb}, err
}

func setupPolygonFinder(dataPath string) (tzf.F, VisualizableTimezoneData, error) {
	var err error
	tzpb := &pb.Timezones{}

	if dataPath == "" {
		internal.Loggger.Debug("Finder will use data from tzf-rel")
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
		internal.Loggger.Debug("Finder will use data from tzf-rel")
		rawFile, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, nil, err
		}
		err = proto.Unmarshal(rawFile, tzpb)
		if err != nil {
			return nil, nil, err
		}
	}
	finder, err := tzf.NewFinderFromPB(tzpb)
	return finder, &polygonData{data: tzpb}, err
}

func NewFinder(cfg *config.Config) (*TZfinder, error) {
	var (
		finder              tzf.F
		err                 error
		tzData              VisualizableTimezoneData
		tzName2Abbreviation = make(map[string]string)
		tzName2Offset       = make(map[string]*durationpb.Duration)
	)

	switch FinderType(cfg.FinderType) {
	case FuzzyFinder:
		finder, tzData, err = setupFuzzyFinder(cfg.DataPath)
	default:
		finder, tzData, err = setupPolygonFinder(cfg.DataPath)
	}

	if err != nil {
		return nil, err
	}

	for _, tzName := range finder.TimezoneNames() {
		location, err := time.LoadLocation(tzName)

		if err != nil {
			return nil, err
		}
		abbreviation, offset := time.Now().In(location).Zone()
		tzName2Abbreviation[tzName] = abbreviation
		tzName2Offset[tzName] = durationpb.New(time.Duration(offset * int(time.Second)))
	}

	return &TZfinder{
		Finder:              finder,
		TZData:              tzData,
		TZName2Abbreviation: tzName2Abbreviation,
		TZName2Offset:       tzName2Offset,
	}, nil

}

var ProviderSet = wire.NewSet(
	NewFinder,
	wire.Bind(new(tzf.F), new(*TZfinder)),
)

func (f *TZfinder) GetTimezoneName(lng, lat float64) string {
	return f.Finder.GetTimezoneName(lng, lat)
}

func (f *TZfinder) GetTimezoneNames(lng, lat float64) ([]string, error) {
	return f.Finder.GetTimezoneNames(lng, lat)
}

func (f *TZfinder) TimezoneNames() []string {
	return f.Finder.TimezoneNames()
}

func (f *TZfinder) DataVersion() string {
	return f.Finder.DataVersion()
}
