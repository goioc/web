/*
 * Copyright (c) 2020 Go IoC
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 */

package web

import (
	"bytes"
	"context"
	"github.com/goioc/di"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var server *httptest.Server

type endpoint1 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint1"`
}

func (e endpoint1) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint1) REST(w http.ResponseWriter) {
	_, err := w.Write([]byte("test"))
	if err != nil {
		panic(err)
	}
}

type endpoint2 struct {
	method interface{} `web.methods:"post,patch"`
	path   interface{} `web.path:"/endpoint2"`
}

func (e endpoint2) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint2) REST(w http.ResponseWriter, r *http.Request) {
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(all)
	if err != nil {
		panic(err)
	}
}

type endpoint3 struct {
	method  interface{} `web.methods:"GET"`
	path    interface{} `web.path:"/endpoint3"`
	queries interface{} `web.queries:"foo,bar,id,{id:[0-9]+}"`
}

func (e endpoint3) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint3) REST(queryParams url.Values) string {
	foo := queryParams.Get("foo")
	id := queryParams.Get("id")
	return foo + id
}

type endpoint4 struct {
	method  interface{} `web.methods:"POST"`
	path    interface{} `web.path:"/endpoint4"`
	headers interface{} `web.headers:"Content-Type,text/plain"`
}

func (e endpoint4) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint4) REST(body string) string {
	return body
}

type endpoint5 struct {
	method  interface{} `web.methods:"GET"`
	path    interface{} `web.path:"/endpoint5/{key}/{*?}"`
	headers interface{} `web.matcher:"matcher"`
}

func (e endpoint5) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint5) REST(pathParams map[string]string) string {
	return pathParams["key"]
}

type endpoint6 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint6"`
}

func (e endpoint6) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint6) REST(ctx context.Context) string {
	return ctx.Value(di.BeanKey("key")).(string)
}

type endpoint7 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint7"`
}

func (e endpoint7) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint7) REST(header http.Header) string {
	return header.Get("Content-Type")
}

type endpoint8 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint8"`
}

func (e endpoint8) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint8) REST(reader io.Reader) string {
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(all)
}

type endpoint9 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint9"`
}

func (e endpoint9) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint9) REST(readCloser io.ReadCloser) string {
	all, err := ioutil.ReadAll(readCloser)
	if err != nil {
		panic(err)
	}
	return string(all)
}

type endpoint10 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint10"`
}

func (e endpoint10) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint10) REST(body []byte) string {
	return string(body)
}

type binaryStruct struct {
	a string
}

func (b *binaryStruct) MarshalBinary() (data []byte, err error) {
	return []byte(b.a), nil
}

func (b *binaryStruct) UnmarshalBinary(data []byte) error {
	b.a = string(data)
	return nil
}

type endpoint11 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint11"`
}

func (e endpoint11) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint11) REST(body binaryStruct) *binaryStruct {
	return &body
}

type textStruct struct {
	a string
}

func (t *textStruct) MarshalBinary() (data []byte, err error) {
	return []byte(t.a), nil
}

func (t *textStruct) UnmarshalText(data []byte) error {
	t.a = string(data)
	return nil
}

type endpoint12 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint12"`
}

func (e endpoint12) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint12) REST(body textStruct) *textStruct {
	return &body
}

type endpoint13 struct {
	method interface{} `web.methods:"POST"`
	path   interface{} `web.path:"/endpoint13"`
}

func (e endpoint13) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint13) REST(body outerStruct) outerStruct {
	return body
}

type endpoint14 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint14"`
}

func (e endpoint14) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint14) REST() (http.Header, int) {
	return map[string][]string{
		"my-header": {"my-header-value"},
	}, 418
}

type endpoint15 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint15"`
}

func (e endpoint15) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint15) REST() []byte {
	return []byte("test")
}

type endpoint16 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint16"`
}

func (e endpoint16) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint16) REST() io.Reader {
	return bytes.NewBufferString("test")
}

type endpoint17 struct {
	method interface{} `web.methods:"GET"`
	path   interface{} `web.path:"/endpoint17"`
}

func (e endpoint17) HandlerFuncName() string {
	return "REST"
}

func (e *endpoint17) REST() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString("test"))
}

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
	_, err := di.RegisterBeanFactory("matcher", di.Singleton, func() (interface{}, error) {
		matcherFunc := mux.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
			return strings.HasSuffix(request.URL.Path, "bar")
		})
		return &matcherFunc, nil
	})
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint1", reflect.TypeOf((*endpoint1)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint2", reflect.TypeOf((*endpoint2)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint3", reflect.TypeOf((*endpoint3)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint4", reflect.TypeOf((*endpoint4)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint5", reflect.TypeOf((*endpoint5)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint6", reflect.TypeOf((*endpoint6)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint7", reflect.TypeOf((*endpoint7)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint8", reflect.TypeOf((*endpoint8)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint9", reflect.TypeOf((*endpoint9)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint10", reflect.TypeOf((*endpoint10)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint11", reflect.TypeOf((*endpoint11)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint12", reflect.TypeOf((*endpoint12)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint13", reflect.TypeOf((*endpoint13)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint14", reflect.TypeOf((*endpoint14)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint15", reflect.TypeOf((*endpoint15)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint16", reflect.TypeOf((*endpoint16)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint17", reflect.TypeOf((*endpoint17)(nil)))
	assert.NoError(suite.T(), err)
	err = di.InitializeContainer()
	assert.NoError(suite.T(), err)
	Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), di.BeanKey("key"), "value")))
		})
	})
	router, err := CreateRouter()
	assert.NoError(suite.T(), err)
	server = httptest.NewServer(router)
}

func (suite *TestSuite) TearDownSuite() {
	server.Close()
}

func TestWebTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestEndpoint1() {
	response, err := http.Get(server.URL + "/endpoint1")
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint2() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint2", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
	buf = bytes.NewBufferString("test")
	request, err := http.NewRequest(http.MethodPatch, server.URL+"/endpoint2", buf)
	assert.NotNil(suite.T(), request)
	assert.NoError(suite.T(), err)
	client := http.Client{}
	assert.NotNil(suite.T(), client)
	response, err = client.Do(request)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err = ioutil.ReadAll(response.Body)
	assert.NotNil(suite.T(), all)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint3() {
	response, err := http.Get(server.URL + "/endpoint3?foo=test&id=42")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, response.StatusCode)
	response, err = http.Get(server.URL + "/endpoint3?foo=bar&id=42")
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "bar42", string(all))
}

func (suite *TestSuite) TestEndpoint4() {
	response, err := http.Post(server.URL+"/endpoint4", "", nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, response.StatusCode)
	buf := bytes.NewBufferString("test")
	response, err = http.Post(server.URL+"/endpoint4", "text/plain", buf)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint5() {
	response, err := http.Get(server.URL + "/endpoint5/foo/baz")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, response.StatusCode)
	response, err = http.Get(server.URL + "/endpoint5/foo/bar")
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "foo", string(all))
}

func (suite *TestSuite) TestEndpoint6() {
	response, err := http.Get(server.URL + "/endpoint6")
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "value", string(all))
}

func (suite *TestSuite) TestEndpoint7() {
	response, err := http.Post(server.URL+"/endpoint7", "application/json", nil)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "application/json", string(all))
}

func (suite *TestSuite) TestEndpoint8() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint8", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint9() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint9", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint10() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint10", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint11() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint11", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint12() {
	buf := bytes.NewBufferString("test")
	response, err := http.Post(server.URL+"/endpoint12", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint13() {
	buf := bytes.NewBufferString(jsonData)
	response, err := http.Post(server.URL+"/endpoint13", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), jsonData, string(all))
}

func (suite *TestSuite) TestEndpoint14() {
	response, err := http.Get(server.URL + "/endpoint14")
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 418, response.StatusCode)
	assert.Equal(suite.T(), "my-header-value", response.Header.Get("my-header"))
}

func (suite *TestSuite) TestEndpoint15() {
	response, err := http.Get(server.URL + "/endpoint15")
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint16() {
	response, err := http.Get(server.URL + "/endpoint16")
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint17() {
	response, err := http.Get(server.URL + "/endpoint17")
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}
