package main

import (
	"fmt"
	"testing"
)

func TestRandomSlice(t *testing.T) {
	slc := []int{1, 2, 3, 4, 5, 6}
	newSlc := RandomSlice(slc)
	for _, item := range newSlc {
		if i, ok := item.(int); ok {
			fmt.Println(i)
		}
	}
}
