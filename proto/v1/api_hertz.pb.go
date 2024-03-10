// Code generated by protoc-gen-go-hertz. DO NOT EDIT.
// versions:
// - protoc-gen-go-hertz v0.3.0
// - protoc             (unknown)
// source: proto/v1/api.proto

package v1

import (
	context "context"
	app "github.com/cloudwego/hertz/pkg/app"
	server "github.com/cloudwego/hertz/pkg/app/server"
	xhertz "github.com/ringsaturn/protoc-gen-go-hertz/xhertz"
	http "net/http"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the Hertz package it is being compiled against.
var _ = new(context.Context)
var _ = new(server.Hertz)
var _ = new(app.Handler)
var _ = new(http.Handler)
var _ = new(xhertz.Error)

const OperationTZFServiceGetAllTimezones = "/proto.v1.TZFService/GetAllTimezones"
const OperationTZFServiceGetTimezone = "/proto.v1.TZFService/GetTimezone"
const OperationTZFServiceGetTimezones = "/proto.v1.TZFService/GetTimezones"

type TZFServiceHTTPServer interface {
	GetAllTimezones(context.Context, *GetAllTimezonesRequest) (*GetAllTimezonesResponse, error)
	GetTimezone(context.Context, *GetTimezoneRequest) (*GetTimezoneResponse, error)
	GetTimezones(context.Context, *GetTimezonesRequest) (*GetTimezonesResponse, error)
}

func RegisterTZFServiceHTTPServer(h *server.Hertz, srv TZFServiceHTTPServer) {
	group := h.Group("/")
	group.Handle("GET", "/api/v1/tz", _TZFService_GetTimezone0_HTTP_Handler(srv))
	group.Handle("GET", "/api/v1/tzs", _TZFService_GetTimezones0_HTTP_Handler(srv))
	group.Handle("GET", "/api/v1/tzs/all", _TZFService_GetAllTimezones0_HTTP_Handler(srv))
}

func _TZFService_GetTimezone0_HTTP_Handler(srv TZFServiceHTTPServer) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		var in GetTimezoneRequest
		if err := ctx.BindAndValidate(&in); err != nil {
			xhertz.HandleBadRequest(ctx, err)
			return
		}

		out, err := srv.GetTimezone(c, &in)
		if err != nil {
			xhertz.HandleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, out)
	}
}

func _TZFService_GetTimezones0_HTTP_Handler(srv TZFServiceHTTPServer) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		var in GetTimezonesRequest
		if err := ctx.BindAndValidate(&in); err != nil {
			xhertz.HandleBadRequest(ctx, err)
			return
		}

		out, err := srv.GetTimezones(c, &in)
		if err != nil {
			xhertz.HandleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, out)
	}
}

func _TZFService_GetAllTimezones0_HTTP_Handler(srv TZFServiceHTTPServer) func(c context.Context, ctx *app.RequestContext) {
	return func(c context.Context, ctx *app.RequestContext) {
		var in GetAllTimezonesRequest
		if err := ctx.BindAndValidate(&in); err != nil {
			xhertz.HandleBadRequest(ctx, err)
			return
		}

		out, err := srv.GetAllTimezones(c, &in)
		if err != nil {
			xhertz.HandleError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, out)
	}
}
