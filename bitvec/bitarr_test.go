package bitvec_test

import (
	"bit-vector/bitvec"
	"fmt"
	"testing"
)

func TestBitArr(t *testing.T) {
	bitArr, err := bitvec.NewBitArr("00000")
	if err != nil{
		t.Error(err)
	}
	bitArr.Set1(0)
	bitArr.Set1(1)
	bitArr.Set0(0)
	fmt.Println(bitArr)
}
