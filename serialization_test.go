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
	"github.com/goioc/di"
	"github.com/stretchr/testify/assert"
)

const jsonData = "{\"A\":\"a\",\"B\":42,\"InnerStruct\":{\"C\":\"42\"}}"

type outerStruct struct {
	A           string
	B           int
	InnerStruct struct {
		C string
	}
}

func (suite *TestSuite) TestSerialize() {
	instance, err := di.GetInstanceSafe(GoiocSerializer)
	assert.NotNil(suite.T(), instance)
	assert.NoError(suite.T(), err)
	serializer := instance.(Serializer)
	object := outerStruct{
		A:           "a",
		B:           42,
		InnerStruct: struct{ C string }{C: "42"},
	}
	serializedBytes, err := serializer.Serialize(object)
	assert.NoError(suite.T(), err)
	serializedString := string(serializedBytes)
	assert.Equal(suite.T(), jsonData, serializedString)
}

func (suite *TestSuite) TestDeserialize() {
	instance, err := di.GetInstanceSafe(GoiocSerializer)
	assert.NotNil(suite.T(), instance)
	assert.NoError(suite.T(), err)
	serializer := instance.(Serializer)
	object := new(outerStruct)
	err = serializer.Deserialize([]byte(jsonData), object)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "a", object.A)
	assert.Equal(suite.T(), 42, object.B)
	assert.Equal(suite.T(), "42", object.InnerStruct.C)
}
