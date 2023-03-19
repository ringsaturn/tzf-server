package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
