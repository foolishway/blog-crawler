package utils

import (
	"fmt"
	"testing"
)

func TestRandomSlice(t *testing.T) {
	slc := []int{1, 2, 3, 4, 5, 6}
	newSlc := RandomSlice(slc)
	fmt.Println(newSlc)
}
