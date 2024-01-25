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
	"encoding/json"
	"github.com/goioc/di"
	"reflect"
)

// GoiocSerializer is an ID for Serializer bean. By default, points to JsonSerializer, but can be overwritten.
const GoiocSerializer = "goiocSerializer"

func init() {
	if _, err := di.RegisterBean(GoiocSerializer, reflect.TypeOf((*JsonSerializer)(nil))); err != nil {
		panic(err)
	}
}

// Serializer interface is used by web library to serialize/deserialize objects. Default implementation: JsonSerializer.
type Serializer interface {
	// Serialize method serializes object to byte array.
	Serialize(interface{}) ([]byte, error)
	// Deserialize method deserializes object from byte array.
	Deserialize([]byte, interface{}) error
}

// JsonSerializer is a default implementation of Serializer interface.
type JsonSerializer struct {
}

// Serialize method serializes object to JSON.
func (js JsonSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize method deserializes object from JSON.
func (js JsonSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
