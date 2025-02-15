package httpserver

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	v1 "github.com/ringsaturn/tzf-server/gen/go/tzf_server/v1"
	"github.com/ringsaturn/tzf-server/internal/finder"
)

//go:embed template/*
var f embed.FS

type GetTimezoneInfoRequest struct {
	Name string  `query:"name"`
	Lng  float64 `query:"longitude"`
	Lat  float64 `query:"latitude"`
}

type GetTimezoneInfoPageResponseItem struct {
	Name string
	URL  string
}

type GetTimezoneInfoPageResponse struct {
	Title string
	Items []*GetTimezoneInfoPageResponseItem
}

type WebHandler struct {
	tzfinder *finder.TZfinder
}

func NewWebHandler(tzfinder *finder.TZfinder) *WebHandler {
	return &WebHandler{
		tzfinder: tzfinder,
	}
}

func (w *WebHandler) Index(c context.Context, ctx *app.RequestContext) {
	ctx.Redirect(http.StatusTemporaryRedirect, []byte("/web/tzs/all"))
}

func (w *WebHandler) GetTimezoneShape(c context.Context, ctx *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	err := ctx.BindAndValidate(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"err": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = w.tzfinder.GetTimezoneName(req.Lng, req.Lat)
	}
	shape, err := w.tzfinder.TZData.GetShape(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, shape)
}

func (w *WebHandler) GetTimezoneInfoPage(c context.Context, ctx *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	if err := ctx.BindAndValidate(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": string(ctx.Request.RequestURI())})
		return
	}

	timezone := w.tzfinder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		ctx.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezone for longitude=%.4f, latitude=%.4f", req.Lng, req.Lat),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
		Name: timezone,
		URL:  fmt.Sprintf("/web/tz/geojson/viewer?name=%v", timezone),
	})

	ctx.HTML(200, "info_multi.html", resp)
}

func (w *WebHandler) GetTimezonesInfoPage(c context.Context, ctx *app.RequestContext) {
	req := &v1.GetTimezoneRequest{}
	err := ctx.BindAndValidate(req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": string(ctx.Request.RequestURI())})
		return
	}

	names, err := w.tzfinder.GetTimezoneNames(req.Longitude, req.Latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezones for longitude=%.4f, latitude=%.4f", req.Longitude, req.Latitude),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	for _, name := range names {
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  fmt.Sprintf("/web/tz/geojson/viewer?name=%v", name),
		})
	}
	ctx.HTML(200, "info_multi.html", resp)
}

func (w *WebHandler) GetAllSupportTimezoneNamesPage(c context.Context, ctx *app.RequestContext) {
	resp := &GetTimezoneInfoPageResponse{
		Title: "All timezones",
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
		Name: "All",
		URL:  fmt.Sprintf("/web/tz/geojson/viewer?name=%v", "all"),
	})

	for _, name := range w.tzfinder.TimezoneNames() {
		viewerURL := fmt.Sprintf("/web/tz/geojson/viewer?name=%v", name)
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  viewerURL,
		})
	}
	ctx.HTML(200, "info_multi.html", resp)
}

func (w *WebHandler) GetGeoJSONViewerForTimezone(c context.Context, ctx *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	if err := ctx.BindAndValidate(req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": string(ctx.Request.RequestURI())})
		return
	}

	timezone := w.tzfinder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		ctx.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	ctx.HTML(http.StatusOK, "viewer.html", map[string]any{
		"URL": fmt.Sprintf("/api/v1/tz/geojson?longitude=%v&latitude=%v&name=%v", req.Lng, req.Lat, req.Name),
	})
}

func (w *WebHandler) GetClickPage(c context.Context, ctx *app.RequestContext) {
	ctx.HTML(http.StatusOK, "click.html", nil)
}

func BindWebPage(h *server.Hertz, w *WebHandler) {
	h.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "template/*.html")))

	h.GET("/", w.Index)
	h.GET("/api/v1/tz/geojson", w.GetTimezoneShape)

	webPageGroup := h.Group("/web")
	webPageGroup.GET("/tz", w.GetTimezoneInfoPage)
	webPageGroup.GET("/tzs", w.GetTimezonesInfoPage)
	webPageGroup.GET("/tzs/all", w.GetAllSupportTimezoneNamesPage)
	webPageGroup.GET("/tz/geojson/viewer", w.GetGeoJSONViewerForTimezone)
	webPageGroup.GET("/click", w.GetClickPage)
}
