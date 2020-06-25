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
