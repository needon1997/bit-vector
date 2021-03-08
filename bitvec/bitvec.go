package bitvec

import (
	"math"
)

type BitVector interface {
	Rank0(index int) int
	Rank1(index int) int
	Select0(index int) int
	Select1(index int) int
}

type BasicBitVector struct {
	bits         *BitArr
	length       int
	blockSize    int
	subBlockSize int
	blocks       []*Block
}
type Block struct {
	rankValue uint
	size      int
	rank1     *BitArr
	blocks    []*Block
}

func (this *BasicBitVector) Rank0(index int) int {
	return index + 1 - this.Rank1(index)
}
func (this *BasicBitVector) Rank1(index int) int {
	b := index / this.blockSize
	j := index % this.blockSize
	c := j / this.subBlockSize
	k := j % this.subBlockSize
	var (
		r1, r2, r3 int
	)
	if b > 0 {
		r1 = int(Value(*this.blocks[b-1].rank1))
	} else {
		r1 = 0
	}
	if c > 0 {
		r2 = int(Value(*this.blocks[b].blocks[c-1].rank1))
	} else {
		r2 = 0
	}
	r3 = int(Value(*this.blocks[b].blocks[c].blocks[k].rank1))
	return r1 + r2 + r3
}
func NewBasicBitVec(bitString string) (*BasicBitVector, error) {
	bitVec, err := initBitVecStructure(bitString)
	if err != nil {
		return nil, err
	}
	size := 0
	psize := -1
	for i := 0; i < len(bitVec.blocks); i++ {
		psize = size
		size += bitVec.blocks[i].size
		rankBlockBitArr := ToBitArr(uint(bitVec.bits.Rank1(size - 1)))
		bitVec.blocks[i].rank1 = &rankBlockBitArr
		bitVec.blocks[i].rankValue = Value(rankBlockBitArr)
		subBlockSize := 0
		pSubBlockSize := -1
		for j := 0; j < len(bitVec.blocks[i].blocks); j++ {
			pSubBlockSize = subBlockSize
			subBlockSize += bitVec.blocks[i].blocks[j].size
			rankSubBlockBitArr := ToBitArr(uint(bitVec.bits.Rank1(psize+subBlockSize-1) - bitVec.bits.Rank1(psize-1)))
			bitVec.blocks[i].blocks[j].rank1 = &rankSubBlockBitArr
			bitVec.blocks[i].blocks[j].rankValue = Value(rankSubBlockBitArr)
			for k := 0; k < len(bitVec.blocks[i].blocks[j].blocks); k++ {
				rankBitArr := ToBitArr(uint(bitVec.bits.Rank1(psize+pSubBlockSize+k) - bitVec.bits.Rank1(psize+pSubBlockSize-1)))
				bitVec.blocks[i].blocks[j].blocks[k].rank1 = &rankBitArr
				bitVec.blocks[i].blocks[j].blocks[k].rankValue = Value(rankBitArr)
			}
		}
	}
	return bitVec, err
}

func initBitVecStructure(bitString string) (*BasicBitVector, error) {
	bitArr, err := NewBitArr(bitString)
	if err != nil {
		return nil, err
	}
	strSize := len(bitString)
	bv := &BasicBitVector{bits: bitArr, length: strSize}
	blockSize := int(math.Ceil(math.Log2(float64(strSize)) * math.Log2(float64(len(bitString)))))
	bv.blockSize = blockSize
	blockNum := int(math.Ceil(float64(len(bitString)) / float64(blockSize)))
	subBlockSize := int(math.Ceil(0.5 * math.Log2(float64(strSize))))
	bv.subBlockSize = subBlockSize
	subBlockNum := int(math.Ceil(float64(blockSize) / float64(subBlockSize)))
	bv.blocks = make([]*Block, blockNum)
	for i := 0; i < blockNum-1; i++ {
		bv.blocks[i] = &Block{size: blockSize}
		bv.blocks[i].blocks = make([]*Block, subBlockNum)
		for j := 0; j < subBlockNum-1; j++ {
			bv.blocks[i].blocks[j] = &Block{size: subBlockSize}
			bv.blocks[i].blocks[j].blocks = make([]*Block, subBlockSize)
			for k := 0; k < subBlockSize; k++ {
				bv.blocks[i].blocks[j].blocks[k] = &Block{size: 1}
			}
		}
		bv.blocks[i].blocks[subBlockNum-1] = &Block{size: blockSize - (subBlockNum-1)*subBlockSize}
		bv.blocks[i].blocks[subBlockNum-1].blocks = make([]*Block, blockSize-(subBlockNum-1)*subBlockSize)
		for k := 0; k < blockSize-(subBlockNum-1)*subBlockSize; k++ {
			bv.blocks[i].blocks[subBlockNum-1].blocks[k] = &Block{size: 1}
		}
	}
	bv.blocks[blockNum-1] = &Block{size: strSize - (blockNum-1)*blockSize}
	lSubBlockNum := int(math.Ceil(float64(bv.blocks[blockNum-1].size) / float64(subBlockSize)))
	bv.blocks[blockNum-1].blocks = make([]*Block, lSubBlockNum)
	for j := 0; j < lSubBlockNum-1; j++ {
		bv.blocks[blockNum-1].blocks[j] = &Block{size: subBlockSize}
		bv.blocks[blockNum-1].blocks[j].blocks = make([]*Block, subBlockSize)
		for k := 0; k < subBlockSize; k++ {
			bv.blocks[blockNum-1].blocks[j].blocks[k] = &Block{size: 1}
		}
	}
	bv.blocks[blockNum-1].blocks[lSubBlockNum-1] = &Block{size: bv.blocks[blockNum-1].size - (lSubBlockNum-1)*subBlockSize}
	bv.blocks[blockNum-1].blocks[lSubBlockNum-1].blocks = make([]*Block, bv.blocks[blockNum-1].size-(lSubBlockNum-1)*subBlockSize)
	for k := 0; k < bv.blocks[blockNum-1].blocks[lSubBlockNum-1].size; k++ {
		bv.blocks[blockNum-1].blocks[lSubBlockNum-1].blocks[k] = &Block{size: 1}
	}
	return bv, nil
}
