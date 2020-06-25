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
	"io/ioutil"
	"net/http"
)

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
