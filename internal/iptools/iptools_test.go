package iptools

import (
	"fmt"
	"testing"
)

func TestIP4Re(t *testing.T) {
	str := "aaa1.1.11.23a36.36.32.200"
	fmt.Println(GetIP4FromString(str))
	fmt.Println(ValidIP4(str))
}
