package bitvec

import (
	"errors"
	"math"
)

type BitArr struct {
	arr    []uint8
	length int
}

func NewBitArr(bitString string) (*BitArr, error) {
	l := len(bitString)
	size := int(math.Ceil(float64(l) / float64(8)))
	barr := BitArr{arr: make([]uint8, size), length: l}
	for i := 0; i < l; i++ {
		if bitString[i] == 48 {
			barr.Set0(i)
		} else if bitString[i] == 49 {
			barr.Set1(i)
		} else {
			return nil, errors.New("wrong bit string format")
		}
	}
	return &barr, nil
}

func Value(barr BitArr) uint {
	var val uint = 1
	l := len(barr.arr)
	var sum uint = 0
	for i := l - 1; i >= 0; i-- {
		sum += uint(barr.arr[i]) * val
		val *= val << 8
	}
	return sum
}
func ToBitArr(val uint) BitArr {
	if val == 0 {
		return BitArr{arr: make([]uint8, 1), length: 8}
	}
	if val < 0 {
		panic("not support")
	}
	cval := val
	blockSize := 0
	for val != 0 {
		val = val >> 7
		blockSize += 1
	}
	arr := make([]uint8, blockSize)
	base := uint(256)
	for i := blockSize - 1; i >= 0; i-- {
		arr[i] = uint8(cval % base)
		cval = cval >> 8
	}
	return BitArr{arr: arr, length: blockSize * 8}
}

func (this *BitArr) Get(i int) uint8 {
	if i >= this.length || i < 0 {
		panic("index out of bound")
	}
	superIndex := i / 8
	index := i % 8
	block := this.arr[superIndex]
	ops := uint8(1) << (7 - index)
	result := block & ops >> (7 - index)
	return result

}
func (this *BitArr) Set1(i int) {
	if i >= this.length || i < 0 {
		panic("index out of bound")
	}
	superIndex := i / 8
	index := i % 8
	ops := uint8(1) << (7 - index)
	this.arr[superIndex] = this.arr[superIndex] | ops
}
func (this *BitArr) Set0(i int) {
	if i >= this.length || i < 0 {
		panic("index out of bound")
	}
	superIndex := i / 8
	index := i % 8
	ops := uint8(uint8(255) - uint8(1)<<(7-index))
	this.arr[superIndex] = this.arr[superIndex] & ops
}

func (this *BitArr) Rank0(index int) int {
	rank := 0
	for i := 0; i <= index; i++ {
		if this.Get(i) == 0 {
			rank += 1
		}
	}
	return rank
}

func (this *BitArr) Rank1(index int) int {
	return index + 1 - this.Rank0(index)
}
