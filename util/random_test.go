package util

import (
	"fmt"
	"testing"
)

func TestApproach8(t *testing.T) {
	fmt.Println(RandNChar(2))
	fmt.Println(RandNChar(2))
	fmt.Println(RandNChar(2))
}
func BenchmarkApproach8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandNChar(10)
	}
}
