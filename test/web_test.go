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

package test

import (
	"bytes"
	"context"
	"github.com/goioc/di"
	"github.com/goioc/web"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var server *httptest.Server

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
	_, err = di.RegisterBean("endpoint1", reflect.TypeOf((*Endpoint1)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint2", reflect.TypeOf((*Endpoint2)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint3", reflect.TypeOf((*Endpoint3)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint4", reflect.TypeOf((*Endpoint4)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint5", reflect.TypeOf((*Endpoint5)(nil)))
	assert.NoError(suite.T(), err)
	_, err = di.RegisterBean("endpoint6", reflect.TypeOf((*Endpoint6)(nil)))
	assert.NoError(suite.T(), err)
	err = di.InitializeContainer()
	assert.NoError(suite.T(), err)
	web.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "key", "value")))
		})
	})
	router, err := web.CreateRouter()
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
	response, err := http.Post(server.URL + "/endpoint2", "", buf)
	assert.NotNil(suite.T(), response)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NotNil(suite.T(), all)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
	buf = bytes.NewBufferString("test")
	request, err := http.NewRequest(http.MethodPatch, server.URL + "/endpoint2", buf)
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
	assert.Equal(suite.T(),  404, response.StatusCode)
	response, err = http.Get(server.URL + "/endpoint3?foo=bar&id=42")
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "bar42", string(all))
}


func (suite *TestSuite) TestEndpoint4() {
	response, err := http.Post(server.URL + "/endpoint4", "", nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(),  404, response.StatusCode)
	buf := bytes.NewBufferString("test")
	response, err = http.Post(server.URL + "/endpoint4", "text/plain", buf)
	assert.NoError(suite.T(), err)
	all, err := ioutil.ReadAll(response.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", string(all))
}

func (suite *TestSuite) TestEndpoint5() {
	response, err := http.Get(server.URL + "/endpoint5/foo/baz")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(),  404, response.StatusCode)
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