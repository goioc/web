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
