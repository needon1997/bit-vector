package bitvec

import (
	"errors"
	"fmt"
	"math"
)

type BitArr struct {
	arr    []uint8
	length int
}

func NewBitArrBySize(n int) *BitArr {
	size := int(math.Ceil(float64(n) / float64(8)))
	barr := BitArr{arr: make([]uint8, size), length: n}
	return &barr
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
		val = val >> 8
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

func (this *BitArr) GetValueInRange(start, end int) uint {
	if end < start {
		panic("end < start")
	} else if start < 0 {
		panic("start < 0")
	} else if end >= this.length {
		panic("index out of bound")
	}
	startSuperIndex := start / 8
	startIndex := start % 8
	endSuperIndex := end / 8
	endIndex := end % 8
	var val uint = 0
	if startSuperIndex != endSuperIndex {
		val += uint(this.arr[endSuperIndex]) >> (8 - (endIndex + 1))
		for i := 0; i < endSuperIndex-startSuperIndex-1; i++ {
			val += uint(this.arr[endSuperIndex-i-1]) << (endIndex + 1 + i*8)
		}
		val += uint(this.arr[startSuperIndex]<<startIndex>>startIndex) << (endIndex + 1 + (endSuperIndex-startSuperIndex-1)*8)
	} else {
		val += uint(this.arr[startSuperIndex] << startIndex >> (startIndex + (8 - (endIndex + 1))))
	}
	return val
}
func (this *BitArr) MapValueBounded(start, end int, val uint) {
	if end < start {
		panic("end < start")
	}
	ba := ToBitArr(val)
	for i := 0; i <= end-start; i++ {
		if ba.length-1-i < 0 {
			this.Set0(end - i)
			continue
		}
		v := ba.Get(ba.length - 1 - i)
		if v == 1 {
			this.Set1(end - i)
		} else {
			this.Set0(end - i)
		}
	}
}
func (this *BitArr) Get(i int) uint8 {
	if i >= this.length || i < 0 {
		fmt.Println("error")
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
func (this *BitArr) RangeRank0(start int, end int) int {
	rank := 0
	for i := start; i <= end; i++ {
		if this.Get(i) == 0 {
			rank += 1
		}
	}
	return rank
}
func (this *BitArr) RangeRank1(start int, end int) int {
	return end - start + 1 - this.RangeRank0(start, end)
}
func (this *BitArr) Rank1(index int) int {
	return index + 1 - this.Rank0(index)
}
