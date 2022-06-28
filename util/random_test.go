package util

import (
	"fmt"
	"testing"
)

func TestApproach8(t *testing.T) {
	fmt.Println((2))
	fmt.Println(RandNCharAccount(2))
	fmt.Println(RandNCharAccount(2))
}
func BenchmarkApproach8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandNCharAccount(10)
	}
}
