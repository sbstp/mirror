package mirror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilPointerInTarget(t *testing.T) {
	type bax struct {
		G string
	}
	type foo struct {
		X *bax
	}

	target := foo{}
	src := foo{
		X: &bax{
			G: "hello",
		},
	}

	DeepCopyInto(&target).From(src)

	assert.Equal(t, "hello", target.X.G)
}

func TestSliceIgnoreZero(t *testing.T) {
	type bax struct {
		Items []int
	}

	target := bax{
		Items: []int{1, 2, 3},
	}
	DeepCopyInto(&target).
		SetIgnoreZeroValues(true).
		From(bax{Items: nil})

	assert.Equal(t, []int{1, 2, 3}, target.Items)
}

func TestMapIgnoreZero(t *testing.T) {
	type bax struct {
		Items map[int]int
	}

	target := bax{
		Items: map[int]int{1: 2},
	}
	DeepCopyInto(&target).
		SetIgnoreZeroValues(true).
		From(bax{Items: nil})

	assert.Equal(t, map[int]int{1: 2}, target.Items)
}

func TestSliceNilOverride(t *testing.T) {
	type bax struct {
		Items []int
	}

	target := bax{
		Items: []int{1, 2, 3},
	}
	DeepCopyInto(&target).
		From(bax{Items: nil})

	assert.Equal(t, []int(nil), target.Items)
}

func TestMapNilOverride(t *testing.T) {
	type bax struct {
		Items map[int]int
	}

	target := bax{
		Items: map[int]int{1: 2},
	}
	DeepCopyInto(&target).
		From(bax{Items: nil})

	assert.Equal(t, map[int]int(nil), target.Items)
}

func TestCopyWithinSlices(t *testing.T) {
	type bax struct {
		Bax string
	}
	type foo struct {
		Z []bax
	}

	target := &foo{
		[]bax{{Bax: "hello"}},
	}
	DeepCopyInto(target).SetIgnoreZeroValues(true).From(foo{
		Z: []bax{{Bax: ""}},
	})

	assert.Equal(t, "hello", target.Z[0].Bax)
}

func TestCopyWithinMaps(t *testing.T) {
	type bax struct {
		A string
		B string
	}

	target := map[string]bax{
		"hello": {
			A: "hello",
		},
	}

	source := map[string]bax{
		"hello": {
			B: "world",
		},
	}

	DeepCopyInto(&target).SetIgnoreZeroValues(true).From(&source)

	result := target["hello"]
	assert.Equal(t, "hello", result.A)
	assert.Equal(t, "world", result.B)
}
