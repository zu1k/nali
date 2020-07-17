package qqwry

import (
	"fmt"
	"testing"
)

func TestNewQQwry(t *testing.T) {
	fmt.Println("看看中文")
	qqwry := NewQQwry("../../db/qqwry.dat")
	fmt.Println(qqwry.Find("8.8.8.8"))
	fmt.Println("我是中文")
}
