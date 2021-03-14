package bitvec

import (
	"math"
	"unsafe"
)

type BitVector interface {
	Rank0(index int) int
	Rank1(index int) int
	Select0(index int) int
	Select1(index int) int
}

type BasicBitVector struct {
	length                      int
	blockSize                   int
	subBlockSize                int
	blockRankBitsNum            int
	subBlockRankBitsNum         int
	subBlockBitRankBitsNum      int
	subBlockBitRankIndexBitsNum int
	subBlockNum                 int
	superBlockRank              *BitArr
	subBlockRank                *BitArr
	subBlockBitRank             *BitArr
	subBlockBitRankIndex        *BitArr
}

func (this *BasicBitVector) SizeOf() uintptr {
	return unsafe.Sizeof(*this) + unsafe.Sizeof(*this.superBlockRank) + unsafe.Sizeof(*this.subBlockBitRank) + unsafe.Sizeof(*this.subBlockRank) + unsafe.Sizeof(*this.subBlockBitRankIndex)
}
func (this *BasicBitVector) Rank0(index int) int {
	return index + 1 - this.Rank1(index)
}
func (this *BasicBitVector) Rank1(index int) int {
	b := index / this.blockSize
	j := index % this.blockSize
	c := j / this.subBlockSize
	k := j % this.subBlockSize
	var rank1, rank2, rank3Index uint
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
	rank3Index = this.subBlockBitRankIndex.GetValueInRange((b*this.subBlockNum+c)*this.subBlockBitRankIndexBitsNum, (b*this.subBlockNum+c+1)*this.subBlockBitRankIndexBitsNum-1)
	rank3 := this.subBlockBitRank.GetValueInRange((int(rank3Index)*this.subBlockSize+k)*this.subBlockBitRankBitsNum, (int(rank3Index)*this.subBlockSize+k+1)*this.subBlockBitRankBitsNum-1)
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

	sqrtN := int(math.Ceil(math.Pow(2, float64(subBlockSize))))
	subBlockBitRankIndexBitsNum := int(math.Ceil(math.Log2(float64(sqrtN))))

	bv := &BasicBitVector{length: arrSize}
	bv.blockSize = blockSize
	bv.subBlockSize = subBlockSize
	bv.superBlockRank = NewBitArrBySize(blockRankBitsNum * blockNum)
	bv.subBlockRank = NewBitArrBySize(blockNum * subBlockNum * subBlockRankBitsNum)
	bv.subBlockBitRank = NewBitArrBySize(sqrtN * subBlockSize * subBlockBitRankBitsNum)
	bv.subBlockBitRankIndex = NewBitArrBySize(blockNum * subBlockNum * subBlockBitRankIndexBitsNum)
	bv.blockRankBitsNum = blockRankBitsNum
	bv.subBlockRankBitsNum = subBlockRankBitsNum
	bv.subBlockBitRankBitsNum = subBlockBitRankBitsNum
	bv.subBlockBitRankIndexBitsNum = subBlockBitRankIndexBitsNum
	bv.subBlockNum = subBlockNum
	subBlockIndex := 0
	prevBlockRank := 0
	//initialize the sqrt(n)  unique bit string with length log n/2
	for i := 0; i < sqrtN; i++ {
		temp := NewBitArrBySize(subBlockSize)
		temp.MapValueBounded(0, subBlockSize-1, uint(i))
		for j := 0; j < subBlockSize; j++ {
			rank := temp.Rank1(j)
			bv.subBlockBitRank.MapValueBounded((i*subBlockSize+j)*subBlockBitRankBitsNum, (i*subBlockSize+j+1)*subBlockBitRankBitsNum-1, uint(rank))
		}
	}
	for i := 1; i <= blockNum; i++ {
		var blockRankVal int
		if i*blockSize < arrSize {
			blockRankVal = bitArr.RangeRank1((i-1)*blockSize, i*blockSize-1) + prevBlockRank
		} else {
			blockRankVal = bitArr.RangeRank1((i-1)*blockSize, arrSize-1) + prevBlockRank
		}
		bv.superBlockRank.MapValueBounded((i-1)*blockRankBitsNum, i*blockRankBitsNum-1, uint(blockRankVal))
		prevSubBlockRank := 0
		subBlockLoopEnded := false
		for j := 1; j <= subBlockNum && !subBlockLoopEnded; j++ {
			var subBlockRankVal int
			if j*subBlockSize+(i-1)*blockSize < arrSize {
				if j*subBlockSize < blockSize {
					subBlockRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, (i-1)*blockSize+j*subBlockSize-1) + prevSubBlockRank
					val := bitArr.GetValueInRange((j-1)*subBlockSize+(i-1)*blockSize, j*subBlockSize+(i-1)*blockSize-1)
					bv.subBlockBitRankIndex.MapValueBounded(subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j-1), subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j)-1, val)
				} else {
					subBlockRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, i*blockSize-1) + prevSubBlockRank
					leftBits := bitArr.GetValueInRange((j-1)*subBlockSize+(i-1)*blockSize, i*blockSize-1)
					val := leftBits << (j*subBlockSize - blockSize)
					bv.subBlockBitRankIndex.MapValueBounded(subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j-1), subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j)-1, val)
				}
			} else {
				subBlockRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, arrSize-1)
				leftBits := bitArr.GetValueInRange((j-1)*subBlockSize+(i-1)*blockSize, arrSize-1)
				val := leftBits << (j*subBlockSize + (i-1)*blockSize - arrSize)
				bv.subBlockBitRankIndex.MapValueBounded(subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j-1), subBlockBitRankIndexBitsNum*((i-1)*subBlockNum+j)-1, val)
				subBlockLoopEnded = true
			}
			bv.subBlockRank.MapValueBounded(subBlockIndex*subBlockRankBitsNum, (1+subBlockIndex)*subBlockRankBitsNum-1, uint(subBlockRankVal))
			subBlockIndex += 1
			//subBlockBitLoopEnded := false
			//for k := 1; k <= subBlockSize && !subBlockBitLoopEnded; k++ {
			//	var subBlockBitRankVal int
			//	if (i-1)*blockSize+(j-1)*subBlockSize+k < arrSize {
			//		if (j-1)*subBlockSize+k < blockSize {
			//			subBlockBitRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, (i-1)*blockSize+(j-1)*subBlockSize+k-1)
			//		} else {
			//			subBlockBitRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, i*blockSize-1)
			//			subBlockBitLoopEnded = true
			//		}
			//	} else {
			//		subBlockBitRankVal = bitArr.RangeRank1((i-1)*blockSize+(j-1)*subBlockSize, arrSize-1)
			//		subBlockBitLoopEnded = true
			//	}
			//	bv.subBlockBitRank.MapValueBounded(subBlockBitIndex*subBlockBitRankBitsNum, (1+subBlockBitIndex)*subBlockBitRankBitsNum-1, uint(subBlockBitRankVal))
			//	subBlockBitIndex += 1
			//}
			prevSubBlockRank = subBlockRankVal
		}
		prevBlockRank = blockRankVal
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
