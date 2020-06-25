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
	"net/http"
)

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
