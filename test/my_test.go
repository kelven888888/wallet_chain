package test

import (
	"fmt"
	"testing"
	"time"
)

func Add(a, b int) int {
	return a + b
}
func TestAdd(t *testing.T) {

	fmt.Println(Add(1, 2))
}
func TestBbb(t *testing.T) {

	fmt.Println(time.Now())
}
