package hashcash

import (
	"crypto/rand"
	"fmt"
	"github.com/Kirill-Znamenskiy/kztest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func newRandomHC(t *testing.T, bits int) *Hashcash {
	resource := make([]byte, 5)
	_, err := rand.Read(resource)
	require.NoError(t, err)
	ret, err := New(bits, fmt.Sprintf("%x", resource))
	require.NoError(t, err)
	return ret
}
func TestCompute(t *testing.T) {
	hc := newRandomHC(t, 5)
	kztest.RunTests(t, hc.Compute, []kztest.TestKit{
		{Arg1: uint64(math.MaxUint64), Result1: ErrMaxAttemtpsIsTooBig},
		{Arg1: uint64(1), Result1: ErrComputingMaxAttemptsExceeded},
		{Arg1: uint64(2), Result1: ErrComputingMaxAttemptsExceeded},
		{Arg1: uint64(3), Result1: ErrComputingMaxAttemptsExceeded},
		{Arg1: uint64(0), Result1: assert.NoError},
	})
	assert.True(t, hc.IsCorrect())
}

func TestIsCorrect(t *testing.T) {

	newhc := func(bits int) *Hashcash { return newRandomHC(t, bits) }

	hc3computed := newhc(3)
	err := hc3computed.Compute(0)
	require.NoError(t, err)

	kztest.RunTests(t, IsCorrect, []kztest.TestKit{
		{Arg1: "", Arg2: -1, Result1: false},
		{Arg1: "", Arg2: 0, Result1: false},
		{Arg1: "", Arg2: +1, Result1: false},
		{Arg1: "any", Arg2: -1, Result1: false},
		{Arg1: "any", Arg2: 0, Result1: true},
		{Arg1: "any", Arg2: +1, Result1: false},
		{Arg1: newhc(3).String(), Arg2: 3, Result1: false},
		{Arg1: hc3computed.String(), Arg2: 0, Result1: true},
		{Arg1: hc3computed.String(), Arg2: 1, Result1: true},
		{Arg1: hc3computed.String(), Arg2: 2, Result1: true},
		{Arg1: hc3computed.String(), Arg2: 3, Result1: true},
		{Arg1: hc3computed.String(), Arg2: 4, Result1: false},
		{Arg1: hc3computed.String(), Arg2: 5, Result1: false},
	})
}
