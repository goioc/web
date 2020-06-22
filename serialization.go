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
	"encoding/json"
	"github.com/goioc/di"
	"reflect"
)

const ResponseSerializer = "responseSerializer"

func init() {
	if _, err := di.RegisterBean(ResponseSerializer, reflect.TypeOf((*JsonSerializer)(nil))); err != nil {
		panic(err)
	}
}

type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize([]byte, interface{}) error
}

type JsonSerializer struct {
}

func (js JsonSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (js JsonSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
