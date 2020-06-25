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
	assert.NotNil(suite.T(), all)
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
