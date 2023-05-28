package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type GetTimezoneInfoPageResponseItem struct {
	Name string
	URL  string
}

type GetTimezoneInfoPageResponse struct {
	Title string
	Items []*GetTimezoneInfoPageResponseItem
}

func GetTimezoneInfoPage(ctx context.Context, c *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}

	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezone for lng=%.4f, lat=%.4f", req.Lng, req.Lat),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
		Name: timezone,
		URL:  fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host()), timezone),
	})

	c.HTML(200, "info_multi.html", resp)
}

func GetTimezonesInfoPage(ctx context.Context, c *app.RequestContext) {
	req := &LocationRequest{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": string(c.Request.RequestURI())})
		return
	}

	names, err := finder.GetTimezoneNames(req.Lng, req.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezones for lng=%.4f, lat=%.4f", req.Lng, req.Lat),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	for _, name := range names {
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host()), name),
		})
	}
	c.HTML(200, "info_multi.html", resp)
}

func GetAllSupportTimezoneNamesPage(ctx context.Context, c *app.RequestContext) {
	resp := &GetTimezoneInfoPageResponse{
		Title: "All timezones",
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
		Name: "All",
		URL:  fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host()), "all"),
	})

	for _, name := range finder.TimezoneNames() {
		viewerURL := fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host()), name)
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  viewerURL,
		})
	}
	c.HTML(200, "info_multi.html", resp)
}

// Render GeoJSON on web
func GetGeoJSONViewerForTimezone(ctx context.Context, c *app.RequestContext) {
	req := &GetTimezoneInfoRequest{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.H{"err": err.Error(), "uri": string(c.Request.RequestURI())})
		return
	}

	timezone := finder.GetTimezoneName(req.Lng, req.Lat)
	if timezone == "" {
		c.JSON(http.StatusInternalServerError, utils.H{"err": "no timezone found"})
		return
	}

	c.HTML(http.StatusOK, "viewer.html", map[string]any{
		"URL": fmt.Sprintf("http://%v/api/v1/tz/geojson?lng=%v&lat=%v&name=%v", string(c.Request.Host()), req.Lng, req.Lat, req.Name),
	})
}
