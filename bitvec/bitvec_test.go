package bitvec_test

import (
	bitvec2 "bit-vector/bitvec"
	"testing"
)

func TestNewBasicBitVec(t *testing.T) {
	str := "1000110000000000"
	bitvec, _ := bitvec2.NewBasicBitVecFromString(str)
	bitArr, _ := bitvec2.NewBitArr(str)
	for i := 0; i < len(str); i++ {
		r1 := bitvec.Rank1(i)
		r2 := bitArr.Rank1(i)
		if r1 != r2 {
			t.Error("wrong")
		}
	}
}
