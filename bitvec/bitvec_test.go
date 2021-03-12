package bitvec_test

import (
	bitvec2 "bit-vector/bitvec"
	"testing"
)

func TestNewBasicBitVec(t *testing.T) {
	str := "1010010011101010101101111110001100111111101000000010100100101010001101000111001111101101111111000011010100011001110111101000110110001010011101110001110000011101111000110010011011001011111110100101010000111"
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

func BenchmarkNewBasicBitVec(b *testing.B) {
	str := "1010010011101010101101111110001100111111101000000010100100101010001101000111001111101101111111000011010100011001110111101000110110001010011101110001110000011101111000110010011011001011111110100101010000111"
	bitvec, _ := bitvec2.NewBasicBitVecFromString(str)
	for i := 0; i < len(str); i++ {
		bitvec.Rank1(i)
	}
}
