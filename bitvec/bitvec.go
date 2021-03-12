package bitvec

import (
	"fmt"
	"math"
)

type BitVector interface {
	Rank0(index int) int
	Rank1(index int) int
	Select0(index int) int
	Select1(index int) int
}

type BasicBitVector struct {
	bits                   *BitArr
	length                 int
	blockSize              int
	subBlockSize           int
	blockRankBitsNum       int
	subBlockRankBitsNum    int
	subBlockBitRankBitsNum int
	subBlockNum            int
	superBlockRank         *BitArr
	subBlockRank           *BitArr
	subBlockBitRank        *BitArr
}

func (this *BasicBitVector) Rank0(index int) int {
	return index + 1 - this.Rank1(index)
}
func (this *BasicBitVector) Rank1(index int) int {
	b := index / this.blockSize
	j := index % this.blockSize
	c := j / this.subBlockSize
	k := j % this.subBlockSize
	var rank1, rank2, rank3 uint
	if b == 0 {
		rank1 = 0
	} else {
		rank1 = this.superBlockRank.GetValueInRange((b-1)*this.blockRankBitsNum, b*this.blockRankBitsNum-1)
	}
	if c == 0 {
		rank2 = 0
	} else {
		rank2 = this.subBlockRank.GetValueInRange((b*this.subBlockNum+c-1)*this.subBlockRankBitsNum, (b*this.subBlockNum+c)*this.subBlockRankBitsNum-1)
	}
	rank3 = this.subBlockBitRank.GetValueInRange((b*this.blockSize+c*this.subBlockSize+k)*this.subBlockBitRankBitsNum, (b*this.blockSize+c*this.subBlockSize+k+1)*this.subBlockBitRankBitsNum-1)
	return int(rank1 + rank2 + rank3)
}

func NewBasicBitVec(bitArr *BitArr) *BasicBitVector {
	arrSize := bitArr.length

	blockRankBitsNum := int(math.Ceil(math.Log2(float64(arrSize + 1))))
	blockSize := int(math.Ceil(math.Log2(float64(arrSize)) * math.Log2(float64(arrSize))))
	blockNum := int(math.Ceil(float64(arrSize) / float64(blockSize)))

	subBlockRankBitsNum := int(math.Ceil(math.Log2(float64(arrSize + 1))))
	subBlockSize := int(math.Ceil(0.5 * math.Log2(float64(arrSize))))
	subBlockNum := int(math.Ceil(float64(blockSize) / float64(subBlockSize)))
	subBlockBitRankBitsNum := int(math.Ceil(math.Log2(float64(subBlockSize + 1))))

	bv := &BasicBitVector{bits: bitArr, length: arrSize}
	bv.blockSize = blockSize
	bv.subBlockSize = subBlockSize
	bv.superBlockRank = NewBitArrBySize(blockRankBitsNum * blockNum)
	bv.subBlockRank = NewBitArrBySize(blockNum * subBlockNum * subBlockRankBitsNum)
	bv.subBlockBitRank = NewBitArrBySize(blockNum * subBlockNum * subBlockSize * subBlockBitRankBitsNum)
	bv.blockRankBitsNum = blockRankBitsNum
	bv.subBlockRankBitsNum = subBlockRankBitsNum
	bv.subBlockBitRankBitsNum = subBlockBitRankBitsNum
	bv.subBlockNum = subBlockNum
	subBlockIndex := 0
	subBlockBitIndex := 0
	prevBlockRank := 0
	for i := 1; i <= blockNum; i++ {
		var blockRankVal int
		if i*blockSize < arrSize {
			blockRankVal = bv.bits.RangeRank1((i-1)*blockSize, i*blockSize-1) + prevBlockRank
		} else {
			blockRankVal = bv.bits.RangeRank1((i-1)*blockSize, arrSize-1) + prevBlockRank
		}
		bv.superBlockRank.MapValueBounded((i-1)*blockRankBitsNum, i*blockRankBitsNum-1, uint(blockRankVal))
		prevSubBlockRank := 0
		subBlockLoopEnded := false
		for j := 1; j <= subBlockNum && !subBlockLoopEnded; j++ {
			var subBlockRankVal int
			if j*subBlockSize+(i-1)*blockSize < arrSize {
				if j*subBlockSize < blockSize {
					subBlockRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, (i-1)*blockSize+j*subBlockSize-1) + prevSubBlockRank
				} else {
					subBlockRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, i*blockSize-1) + prevSubBlockRank
				}
			} else {
				subBlockRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, arrSize-1)
				subBlockLoopEnded = true
			}
			bv.subBlockRank.MapValueBounded(subBlockIndex*subBlockRankBitsNum, (1+subBlockIndex)*subBlockRankBitsNum-1, uint(subBlockRankVal))
			subBlockIndex += 1
			subBlockBitLoopEnded := false
			for k := 1; k <= subBlockSize && !subBlockBitLoopEnded; k++ {
				var subBlockBitRankVal int
				if (i-1)*blockSize+(j-1)*subBlockSize+k < arrSize {
					if (j-1)*subBlockSize+k < blockSize {
						subBlockBitRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, (i-1)*blockSize+(j-1)*subBlockSize+k-1)
					} else {
						subBlockBitRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, i*blockSize-1)
						subBlockBitLoopEnded = true
					}
				} else {
					subBlockBitRankVal = bv.bits.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, arrSize-1)
					subBlockBitLoopEnded = true
				}
				bv.subBlockBitRank.MapValueBounded(subBlockBitIndex*subBlockBitRankBitsNum, (1+subBlockBitIndex)*subBlockBitRankBitsNum-1, uint(subBlockBitRankVal))
				subBlockBitIndex += 1
			}
			prevSubBlockRank = subBlockRankVal
		}
		prevBlockRank = blockRankVal
		if i%10 == 0 {
			fmt.Println(i)
		}
	}

	return bv
}

func NewBasicBitVecFromString(bitstring string) (*BasicBitVector, error) {
	bitArr, err := NewBitArr(bitstring)
	if err != nil {
		return nil, err
	}
	return NewBasicBitVec(bitArr), nil
}
