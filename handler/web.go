package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetTimezoneInfoPageResponseItem struct {
	Name string
	URL  string
}

type GetTimezoneInfoPageResponse struct {
	Title string
	Items []*GetTimezoneInfoPageResponseItem
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

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezone for lng=%.4f, lat=%.4f", req.Lng, req.Lat),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
		Name: timezone,
		URL:  fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host), timezone),
	})

	c.HTML(200, "info_multi.html", resp)
}

func GetTimezonesInfoPage(c *gin.Context) {
	req := &LocationRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error(), "uri": c.Request.RequestURI})
		return
	}

	names, err := finder.GetTimezoneNames(req.Lng, req.Lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no timezone found"})
		return
	}

	resp := &GetTimezoneInfoPageResponse{
		Title: fmt.Sprintf("Timezones for lng=%.4f, lat=%.4f", req.Lng, req.Lat),
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	for _, name := range names {
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host), name),
		})
	}
	c.HTML(200, "info_multi.html", resp)
}

func GetAllSupportTimezoneNamesPage(c *gin.Context) {
	resp := &GetTimezoneInfoPageResponse{
		Title: "All timezones",
		Items: make([]*GetTimezoneInfoPageResponseItem, 0),
	}

	for _, name := range finder.TimezoneNames() {
		viewerURL := fmt.Sprintf("http://%v/web/tz/geojson/viewer?name=%v", string(c.Request.Host), name)
		resp.Items = append(resp.Items, &GetTimezoneInfoPageResponseItem{
			Name: name,
			URL:  viewerURL,
		})
	}
	c.HTML(200, "info_multi.html", resp)
}

// Render GeoJSON on web
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

	c.HTML(http.StatusOK, "viewer.html", map[string]any{
		"URL": fmt.Sprintf("http://%v/api/v1/tz/geojson?lng=%v&lat=%v&name=%v", c.Request.Host, req.Lng, req.Lat, req.Name),
	})
}
