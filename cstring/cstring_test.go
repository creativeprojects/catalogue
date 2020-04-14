package cstring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateLength(t *testing.T) {
	testData := []struct {
		cstring []int8
		length  int
	}{
		{[]int8{}, 0},
		{[]int8{0}, 0},
		{[]int8{0, 0, 0}, 0},
		{[]int8{100, 105, 115, 107}, 4},
		{[]int8{100, 105, 115, 107, 0}, 4},
		{[]int8{100, 105, 115, 107, 0, 0, 0}, 4},
	}
	for _, item := range testData {
		length := calculateLen(item.cstring)
		assert.Equal(t, item.length, length)
	}
}
func TestInt8ToBytes(t *testing.T) {
	testData := []struct {
		cstring  []int8
		expected []byte
	}{
		{[]int8{}, []byte{}},
		{[]int8{0}, []byte{}},
		{[]int8{0, 0, 0}, []byte{}},
		{[]int8{100, 105, 115, 107}, []byte("disk")},
		{[]int8{100, 105, 115, 107, 0}, []byte("disk")},
		{[]int8{100, 105, 115, 107, 0, 0, 0}, []byte("disk")},
	}
	for _, item := range testData {
		output := Int8ToBytes(item.cstring)
		assert.Exactly(t, item.expected, output)
	}
}
