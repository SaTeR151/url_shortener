package generator

import (
	"github.com/thanhpk/randstr"
)

func Code(n int) string {
	return randstr.Hex(n)
}
