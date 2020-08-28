package mirror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilPointer(t *testing.T) {
	type bax struct {
		g string
	}
	type foo struct {
		X *bax
	}

	target := foo{}
	src := foo{
		X: &bax{
			g: "hello",
		},
	}

	DeepCopyInto(&target).From(src)

	assert.Equal(t, "hello", target.X.g)
}
