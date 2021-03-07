package bitvec

import "math"

type BitVector interface {
	Rank0(index int) int
	Rank1(index int) int
	Select0(index int) int
	Select1(index int) int
}

type BasicBitVector struct {
	bits   *BitArr
	blocks []Block
}
type Block struct {
	size   int
	rank1  *BitArr
	blocks []Block
}

func NewBasicBitVec(bitString string) (*BasicBitVector, error) {
	bitArr, err := NewBitArr(bitString)
	if err != nil {
		return nil, err
	}
	bv := &BasicBitVector{bits: bitArr}
	blockSize := int(math.Ceil(math.Log2(float64(len(bitString))) * math.Log2(float64(len(bitString)))))
	blockNum := int(math.Ceil(float64(len(bitString) / blockSize)))
	bv.blocks = make([]Block, blockNum)

}
