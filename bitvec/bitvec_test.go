package bitvec_test

import (
	bitvec2 "bit-vector/bitvec"
	"testing"
)

func TestNewBasicBitVec(t *testing.T) {
	str := "10001100000000000000111110001110001111000111000111000111000111001110110101110111011101100000001"
	bitvec, _ := bitvec2.NewBasicBitVec(str)
	bitArr, _ := bitvec2.NewBitArr(str)
	for i := 0; i < len(str); i++ {
		r1 := bitvec.Rank1(i)
		r2 := bitArr.Rank1(i)
		if r1 != r2 {
			t.Error("wrong")
		}
	}
}
