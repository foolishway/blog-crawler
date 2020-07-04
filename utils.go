package main

import (
	"math/rand"
	"reflect"
	"time"
)

func RandomSlice(slc interface{}) []interface{} {
	v := reflect.ValueOf(slc)
	if v.Kind() != reflect.Slice {
		panic("The arg of RandomSlice must be slice.")
	}

	var newSlice []interface{}
	newSlice = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		newSlice[i] = v.Index(i).Interface()
	}
	rand.Seed(time.Now().Unix())
	for i := len(newSlice) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		newSlice[i], newSlice[num] = newSlice[num], newSlice[i]
	}
	return newSlice
}
