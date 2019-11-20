package lib_test

import (
	"testing"

	"github.com/changkun/occamy/lib"
)

func TestNewPool(t *testing.T) {
	p := lib.NewPool(10)

	for i := 0; i < 100; i++ {
		ii := p.Next()
		if ii != i {
			t.Fatalf("want %d, got %d", i, ii)
		}
	}

	for i := 0; i < 100; i++ {
		p.Free(i)
	}

	for i := 0; i < 100; i++ {
		ii := p.Next()
		if ii != i {
			t.Fatalf("want %d, got %d", i, ii)
		}
	}
}