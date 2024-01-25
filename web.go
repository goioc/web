/*
 * Copyright (c) 2024 Go IoC
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
	"context"
	"encoding"
	"github.com/goioc/di"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	textTemplate "text/template"
)

const (
	methods = "web.methods"
	path    = "web.path"
	queries = "web.queries"
	headers = "web.headers"
	matcher = "web.matcher"
)

var middlewareFunctionsInternal []mux.MiddlewareFunc

// Endpoint is an interface representing web endpoint.
type Endpoint interface {
	// HandlerFuncName should return a method name that is going to be used to create http handler.
	HandlerFuncName() string
}

// Use function registers middleware.
func Use(middlewareFunctions ...mux.MiddlewareFunc) {
	middlewareFunctionsInternal = middlewareFunctions
}

// ListenAndServe function wraps http.ListenAndServe(...), automatically creating endpoints from registered beans.
func ListenAndServe(addr string) error {
	router, err := CreateRouter()
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, router)
}

// ListenAndServeTLS function wraps http.ListenAndServeTLS(...), automatically creating endpoints from registered beans.
func ListenAndServeTLS(addr, certFile, keyFile string) error {
	router, err := CreateRouter()
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(addr, certFile, keyFile, router)
}

// CreateRouter function creates *mux.Router (which implements http.Handler interface).
func CreateRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	router.Use(di.Middleware)
	router.Use(middlewareFunctionsInternal...)
	err := registerHandlers(router)
	if err != nil {
		return nil, err
	}
	err = walk(router)
	if err != nil {
		return nil, err
	}
	return router, nil
}

func registerHandlers(router *mux.Router) error {
	logrus.Trace("Registering endpoints...")
	endpointType := reflect.TypeOf((*Endpoint)(nil)).Elem()
	for beanID, beanType := range di.GetBeanTypes() {
		if !beanType.Implements(endpointType) || di.GetBeanScopes()[beanID] != di.Singleton {
			continue
		}
		err := registerHandler(router, beanID, beanType.Elem())
		if err != nil {
			return err
		}
	}
	return nil
}

func registerHandler(router *mux.Router, beanID string, beanType reflect.Type) error {
	endpoint, err := di.GetInstanceSafe(beanID)
	if err != nil {
		return err
	}
	route := router.Name(beanID)
	for i := 0; i < beanType.NumField(); i++ {
		field := beanType.Field(i)
		tag := field.Tag
		if value, ok := tag.Lookup(methods); ok {
			methods := strings.Split(value, ",")
			route = route.Methods(methods...)
		}
		if value, ok := tag.Lookup(path); ok {
			route = route.Path(value)
		}
		if value, ok := tag.Lookup(queries); ok {
			queries := strings.Split(value, ",")
			route = route.Queries(queries...)
		}
		if value, ok := tag.Lookup(headers); ok {
			headers := strings.Split(value, ",")
			route = route.Headers(headers...)
		}
		if value, ok := tag.Lookup(matcher); ok {
			instance, err := di.GetInstanceSafe(value)
			if err != nil {
				return err
			}
			matcher := instance.(*mux.MatcherFunc)
			route = route.MatcherFunc(*matcher)
		}
	}
	route.Handler(createHandler(endpoint.(Endpoint)))
	return nil
}

func createHandler(endpoint Endpoint) http.Handler {
	handlerFunc := reflect.ValueOf(endpoint).MethodByName(endpoint.HandlerFuncName())
	handlerFuncType := handlerFunc.Type()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arguments := make([]reflect.Value, 0)
		for i := 0; i < handlerFuncType.NumIn(); i++ {
			argument := handlerFuncType.In(i)
			switch argument {
			case reflect.TypeOf((*context.Context)(nil)).Elem():
				arguments = append(arguments, reflect.ValueOf(r.Context()))
			case reflect.TypeOf((*http.ResponseWriter)(nil)).Elem():
				arguments = append(arguments, reflect.ValueOf(w))
			case reflect.TypeOf((*http.Request)(nil)):
				arguments = append(arguments, reflect.ValueOf(r))
			case reflect.TypeOf((*http.Header)(nil)).Elem():
				arguments = append(arguments, reflect.ValueOf(r.Header))
			case reflect.TypeOf((*io.Reader)(nil)).Elem():
				arguments = append(arguments, reflect.ValueOf(r.Body))
			case reflect.TypeOf((*io.ReadCloser)(nil)).Elem():
				arguments = append(arguments, reflect.ValueOf(r.Body))
			case reflect.TypeOf((*[]byte)(nil)).Elem():
				all, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				arguments = append(arguments, reflect.ValueOf(all))
			case reflect.TypeOf((*string)(nil)).Elem():
				all, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				arguments = append(arguments, reflect.ValueOf(string(all)))
			case reflect.TypeOf((map[string]string)(nil)):
				arguments = append(arguments, reflect.ValueOf(mux.Vars(r)))
			case reflect.TypeOf((url.Values)(nil)):
				arguments = append(arguments, reflect.ValueOf(r.URL.Query()))
			default:
				all, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				body := reflect.New(argument).Interface()
				typeOfBody := reflect.TypeOf(body)
				if typeOfBody.Implements(reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()) {
					binaryUnmarshaler := body.(encoding.BinaryUnmarshaler)
					if err := binaryUnmarshaler.UnmarshalBinary(all); err != nil {
						panic(err)
					}
					arguments = append(arguments, reflect.ValueOf(binaryUnmarshaler).Elem())
				} else if typeOfBody.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
					textUnmarshaler := body.(encoding.TextUnmarshaler)
					if err := textUnmarshaler.UnmarshalText(all); err != nil {
						panic(err)
					}
					arguments = append(arguments, reflect.ValueOf(textUnmarshaler).Elem())
				} else {
					webResponseSerializer, err := di.GetInstanceSafe(GoiocSerializer)
					if err != nil {
						panic(err)
					}
					if err := webResponseSerializer.(Serializer).Deserialize(all, &body); err != nil {
						panic(err)
					}
					arguments = append(arguments, reflect.ValueOf(body).Elem())
				}
			}
		}
		results := handlerFunc.Call(arguments)
	L:
		for i, result := range results {
			value := result.Interface()
			switch result.Type() {
			case reflect.TypeOf((*int)(nil)).Elem():
				w.WriteHeader(value.(int))
			case reflect.TypeOf((*http.Header)(nil)).Elem():
				for k, v := range value.(http.Header) {
					for _, header := range v {
						w.Header().Add(k, header)
					}
				}
			case reflect.TypeOf((*string)(nil)).Elem():
				if _, err := w.Write([]byte(value.(string))); err != nil {
					panic(err)
				}
				break L
			case reflect.TypeOf((*[]byte)(nil)).Elem():
				if _, err := w.Write(value.([]byte)); err != nil {
					panic(err)
				}
				break L
			case reflect.TypeOf((*io.Reader)(nil)).Elem():
				readCloser := value.(io.Reader)
				if _, err := io.Copy(w, readCloser); err != nil {
					panic(err)
				}
				break L
			case reflect.TypeOf((*io.ReadCloser)(nil)).Elem():
				readCloser := value.(io.ReadCloser)
				if _, err := io.Copy(w, readCloser); err != nil {
					panic(err)
				}
				if err := readCloser.Close(); err != nil {
					panic(err)
				}
				break L
			case reflect.TypeOf((*htmlTemplate.Template)(nil)).Elem():
				tmpl := value.(htmlTemplate.Template)
				if err := tmpl.Execute(w, results[i+1].Interface()); err != nil {
					panic(err)
				}
				break L
			case reflect.TypeOf((*textTemplate.Template)(nil)).Elem():
				tmpl := value.(textTemplate.Template)
				if err := tmpl.Execute(w, results[i+1].Interface()); err != nil {
					panic(err)
				}
				break L
			default:
				if result.Type().Implements(reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()) ||
					reflect.PtrTo(result.Type()).Implements(reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()) {
					marshaler := value.(encoding.BinaryMarshaler)
					body, err := marshaler.MarshalBinary()
					if err != nil {
						panic(err)
					}
					if _, err = w.Write(body); err != nil {
						panic(err)
					}
				} else if result.Type().Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) ||
					reflect.PtrTo(result.Type()).Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) {
					marshaler := value.(encoding.TextMarshaler)
					body, err := marshaler.MarshalText()
					if err != nil {
						panic(err)
					}
					if _, err = w.Write(body); err != nil {
						panic(err)
					}
				} else {
					webResponseSerializer, err := di.GetInstanceSafe(GoiocSerializer)
					if err != nil {
						panic(err)
					}
					body, err := webResponseSerializer.(Serializer).Serialize(value)
					if err != nil {
						panic(err)
					}
					if _, err = w.Write(body); err != nil {
						panic(err)
					}
				}
				break L
			}
		}
	})
}

func walk(router *mux.Router) error {
	logrus.Trace("Registered endpoints: ")
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"route":           pathTemplate,
			"methods":         methods,
			"query templates": strings.Join(queriesTemplates, ","),
		}).Trace("Endpoint registered: ", route.GetName())
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
