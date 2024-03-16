// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/ringsaturn/go/pkg/mod/github.com/tidwall/redcon@v1.6.2/redcon.go

// Package handler_test is a generated GoMock package.
package server_test

import (
	net "net"
	reflect "reflect"

	redcon "github.com/tidwall/redcon"
	gomock "go.uber.org/mock/gomock"
)

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// Context mocks base method.
func (m *MockConn) Context() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockConnMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockConn)(nil).Context))
}

// Detach mocks base method.
func (m *MockConn) Detach() redcon.DetachedConn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detach")
	ret0, _ := ret[0].(redcon.DetachedConn)
	return ret0
}

// Detach indicates an expected call of Detach.
func (mr *MockConnMockRecorder) Detach() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detach", reflect.TypeOf((*MockConn)(nil).Detach))
}

// NetConn mocks base method.
func (m *MockConn) NetConn() net.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetConn")
	ret0, _ := ret[0].(net.Conn)
	return ret0
}

// NetConn indicates an expected call of NetConn.
func (mr *MockConnMockRecorder) NetConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetConn", reflect.TypeOf((*MockConn)(nil).NetConn))
}

// PeekPipeline mocks base method.
func (m *MockConn) PeekPipeline() []redcon.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekPipeline")
	ret0, _ := ret[0].([]redcon.Command)
	return ret0
}

// PeekPipeline indicates an expected call of PeekPipeline.
func (mr *MockConnMockRecorder) PeekPipeline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekPipeline", reflect.TypeOf((*MockConn)(nil).PeekPipeline))
}

// ReadPipeline mocks base method.
func (m *MockConn) ReadPipeline() []redcon.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPipeline")
	ret0, _ := ret[0].([]redcon.Command)
	return ret0
}

// ReadPipeline indicates an expected call of ReadPipeline.
func (mr *MockConnMockRecorder) ReadPipeline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPipeline", reflect.TypeOf((*MockConn)(nil).ReadPipeline))
}

// RemoteAddr mocks base method.
func (m *MockConn) RemoteAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockConnMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConn)(nil).RemoteAddr))
}

// SetContext mocks base method.
func (m *MockConn) SetContext(v interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetContext", v)
}

// SetContext indicates an expected call of SetContext.
func (mr *MockConnMockRecorder) SetContext(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContext", reflect.TypeOf((*MockConn)(nil).SetContext), v)
}

// SetReadBuffer mocks base method.
func (m *MockConn) SetReadBuffer(bytes int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReadBuffer", bytes)
}

// SetReadBuffer indicates an expected call of SetReadBuffer.
func (mr *MockConnMockRecorder) SetReadBuffer(bytes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadBuffer", reflect.TypeOf((*MockConn)(nil).SetReadBuffer), bytes)
}

// WriteAny mocks base method.
func (m *MockConn) WriteAny(any interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteAny", any)
}

// WriteAny indicates an expected call of WriteAny.
func (mr *MockConnMockRecorder) WriteAny(any interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAny", reflect.TypeOf((*MockConn)(nil).WriteAny), any)
}

// WriteArray mocks base method.
func (m *MockConn) WriteArray(count int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteArray", count)
}

// WriteArray indicates an expected call of WriteArray.
func (mr *MockConnMockRecorder) WriteArray(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteArray", reflect.TypeOf((*MockConn)(nil).WriteArray), count)
}

// WriteBulk mocks base method.
func (m *MockConn) WriteBulk(bulk []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteBulk", bulk)
}

// WriteBulk indicates an expected call of WriteBulk.
func (mr *MockConnMockRecorder) WriteBulk(bulk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBulk", reflect.TypeOf((*MockConn)(nil).WriteBulk), bulk)
}

// WriteBulkString mocks base method.
func (m *MockConn) WriteBulkString(bulk string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteBulkString", bulk)
}

// WriteBulkString indicates an expected call of WriteBulkString.
func (mr *MockConnMockRecorder) WriteBulkString(bulk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBulkString", reflect.TypeOf((*MockConn)(nil).WriteBulkString), bulk)
}

// WriteError mocks base method.
func (m *MockConn) WriteError(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteError", msg)
}

// WriteError indicates an expected call of WriteError.
func (mr *MockConnMockRecorder) WriteError(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteError", reflect.TypeOf((*MockConn)(nil).WriteError), msg)
}

// WriteInt mocks base method.
func (m *MockConn) WriteInt(num int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteInt", num)
}

// WriteInt indicates an expected call of WriteInt.
func (mr *MockConnMockRecorder) WriteInt(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt", reflect.TypeOf((*MockConn)(nil).WriteInt), num)
}

// WriteInt64 mocks base method.
func (m *MockConn) WriteInt64(num int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteInt64", num)
}

// WriteInt64 indicates an expected call of WriteInt64.
func (mr *MockConnMockRecorder) WriteInt64(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt64", reflect.TypeOf((*MockConn)(nil).WriteInt64), num)
}

// WriteNull mocks base method.
func (m *MockConn) WriteNull() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteNull")
}

// WriteNull indicates an expected call of WriteNull.
func (mr *MockConnMockRecorder) WriteNull() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteNull", reflect.TypeOf((*MockConn)(nil).WriteNull))
}

// WriteRaw mocks base method.
func (m *MockConn) WriteRaw(data []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteRaw", data)
}

// WriteRaw indicates an expected call of WriteRaw.
func (mr *MockConnMockRecorder) WriteRaw(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteRaw", reflect.TypeOf((*MockConn)(nil).WriteRaw), data)
}

// WriteString mocks base method.
func (m *MockConn) WriteString(str string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteString", str)
}

// WriteString indicates an expected call of WriteString.
func (mr *MockConnMockRecorder) WriteString(str interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*MockConn)(nil).WriteString), str)
}

// WriteUint64 mocks base method.
func (m *MockConn) WriteUint64(num uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteUint64", num)
}

// WriteUint64 indicates an expected call of WriteUint64.
func (mr *MockConnMockRecorder) WriteUint64(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteUint64", reflect.TypeOf((*MockConn)(nil).WriteUint64), num)
}

// MockDetachedConn is a mock of DetachedConn interface.
type MockDetachedConn struct {
	ctrl     *gomock.Controller
	recorder *MockDetachedConnMockRecorder
}

// MockDetachedConnMockRecorder is the mock recorder for MockDetachedConn.
type MockDetachedConnMockRecorder struct {
	mock *MockDetachedConn
}

// NewMockDetachedConn creates a new mock instance.
func NewMockDetachedConn(ctrl *gomock.Controller) *MockDetachedConn {
	mock := &MockDetachedConn{ctrl: ctrl}
	mock.recorder = &MockDetachedConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDetachedConn) EXPECT() *MockDetachedConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDetachedConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDetachedConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDetachedConn)(nil).Close))
}

// Context mocks base method.
func (m *MockDetachedConn) Context() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDetachedConnMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDetachedConn)(nil).Context))
}

// Detach mocks base method.
func (m *MockDetachedConn) Detach() redcon.DetachedConn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detach")
	ret0, _ := ret[0].(redcon.DetachedConn)
	return ret0
}

// Detach indicates an expected call of Detach.
func (mr *MockDetachedConnMockRecorder) Detach() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detach", reflect.TypeOf((*MockDetachedConn)(nil).Detach))
}

// Flush mocks base method.
func (m *MockDetachedConn) Flush() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockDetachedConnMockRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockDetachedConn)(nil).Flush))
}

// NetConn mocks base method.
func (m *MockDetachedConn) NetConn() net.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetConn")
	ret0, _ := ret[0].(net.Conn)
	return ret0
}

// NetConn indicates an expected call of NetConn.
func (mr *MockDetachedConnMockRecorder) NetConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetConn", reflect.TypeOf((*MockDetachedConn)(nil).NetConn))
}

// PeekPipeline mocks base method.
func (m *MockDetachedConn) PeekPipeline() []redcon.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekPipeline")
	ret0, _ := ret[0].([]redcon.Command)
	return ret0
}

// PeekPipeline indicates an expected call of PeekPipeline.
func (mr *MockDetachedConnMockRecorder) PeekPipeline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekPipeline", reflect.TypeOf((*MockDetachedConn)(nil).PeekPipeline))
}

// ReadCommand mocks base method.
func (m *MockDetachedConn) ReadCommand() (redcon.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCommand")
	ret0, _ := ret[0].(redcon.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCommand indicates an expected call of ReadCommand.
func (mr *MockDetachedConnMockRecorder) ReadCommand() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCommand", reflect.TypeOf((*MockDetachedConn)(nil).ReadCommand))
}

// ReadPipeline mocks base method.
func (m *MockDetachedConn) ReadPipeline() []redcon.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPipeline")
	ret0, _ := ret[0].([]redcon.Command)
	return ret0
}

// ReadPipeline indicates an expected call of ReadPipeline.
func (mr *MockDetachedConnMockRecorder) ReadPipeline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPipeline", reflect.TypeOf((*MockDetachedConn)(nil).ReadPipeline))
}

// RemoteAddr mocks base method.
func (m *MockDetachedConn) RemoteAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockDetachedConnMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockDetachedConn)(nil).RemoteAddr))
}

// SetContext mocks base method.
func (m *MockDetachedConn) SetContext(v interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetContext", v)
}

// SetContext indicates an expected call of SetContext.
func (mr *MockDetachedConnMockRecorder) SetContext(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContext", reflect.TypeOf((*MockDetachedConn)(nil).SetContext), v)
}

// SetReadBuffer mocks base method.
func (m *MockDetachedConn) SetReadBuffer(bytes int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReadBuffer", bytes)
}

// SetReadBuffer indicates an expected call of SetReadBuffer.
func (mr *MockDetachedConnMockRecorder) SetReadBuffer(bytes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadBuffer", reflect.TypeOf((*MockDetachedConn)(nil).SetReadBuffer), bytes)
}

// WriteAny mocks base method.
func (m *MockDetachedConn) WriteAny(any interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteAny", any)
}

// WriteAny indicates an expected call of WriteAny.
func (mr *MockDetachedConnMockRecorder) WriteAny(any interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAny", reflect.TypeOf((*MockDetachedConn)(nil).WriteAny), any)
}

// WriteArray mocks base method.
func (m *MockDetachedConn) WriteArray(count int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteArray", count)
}

// WriteArray indicates an expected call of WriteArray.
func (mr *MockDetachedConnMockRecorder) WriteArray(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteArray", reflect.TypeOf((*MockDetachedConn)(nil).WriteArray), count)
}

// WriteBulk mocks base method.
func (m *MockDetachedConn) WriteBulk(bulk []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteBulk", bulk)
}

// WriteBulk indicates an expected call of WriteBulk.
func (mr *MockDetachedConnMockRecorder) WriteBulk(bulk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBulk", reflect.TypeOf((*MockDetachedConn)(nil).WriteBulk), bulk)
}

// WriteBulkString mocks base method.
func (m *MockDetachedConn) WriteBulkString(bulk string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteBulkString", bulk)
}

// WriteBulkString indicates an expected call of WriteBulkString.
func (mr *MockDetachedConnMockRecorder) WriteBulkString(bulk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBulkString", reflect.TypeOf((*MockDetachedConn)(nil).WriteBulkString), bulk)
}

// WriteError mocks base method.
func (m *MockDetachedConn) WriteError(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteError", msg)
}

// WriteError indicates an expected call of WriteError.
func (mr *MockDetachedConnMockRecorder) WriteError(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteError", reflect.TypeOf((*MockDetachedConn)(nil).WriteError), msg)
}

// WriteInt mocks base method.
func (m *MockDetachedConn) WriteInt(num int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteInt", num)
}

// WriteInt indicates an expected call of WriteInt.
func (mr *MockDetachedConnMockRecorder) WriteInt(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt", reflect.TypeOf((*MockDetachedConn)(nil).WriteInt), num)
}

// WriteInt64 mocks base method.
func (m *MockDetachedConn) WriteInt64(num int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteInt64", num)
}

// WriteInt64 indicates an expected call of WriteInt64.
func (mr *MockDetachedConnMockRecorder) WriteInt64(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt64", reflect.TypeOf((*MockDetachedConn)(nil).WriteInt64), num)
}

// WriteNull mocks base method.
func (m *MockDetachedConn) WriteNull() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteNull")
}

// WriteNull indicates an expected call of WriteNull.
func (mr *MockDetachedConnMockRecorder) WriteNull() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteNull", reflect.TypeOf((*MockDetachedConn)(nil).WriteNull))
}

// WriteRaw mocks base method.
func (m *MockDetachedConn) WriteRaw(data []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteRaw", data)
}

// WriteRaw indicates an expected call of WriteRaw.
func (mr *MockDetachedConnMockRecorder) WriteRaw(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteRaw", reflect.TypeOf((*MockDetachedConn)(nil).WriteRaw), data)
}

// WriteString mocks base method.
func (m *MockDetachedConn) WriteString(str string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteString", str)
}

// WriteString indicates an expected call of WriteString.
func (mr *MockDetachedConnMockRecorder) WriteString(str interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*MockDetachedConn)(nil).WriteString), str)
}

// WriteUint64 mocks base method.
func (m *MockDetachedConn) WriteUint64(num uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteUint64", num)
}

// WriteUint64 indicates an expected call of WriteUint64.
func (mr *MockDetachedConnMockRecorder) WriteUint64(num interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteUint64", reflect.TypeOf((*MockDetachedConn)(nil).WriteUint64), num)
}

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// ServeRESP mocks base method.
func (m *MockHandler) ServeRESP(conn redcon.Conn, cmd redcon.Command) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeRESP", conn, cmd)
}

// ServeRESP indicates an expected call of ServeRESP.
func (mr *MockHandlerMockRecorder) ServeRESP(conn, cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeRESP", reflect.TypeOf((*MockHandler)(nil).ServeRESP), conn, cmd)
}