package mPOW

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/hashcash"
	"github.com/Kirill-Znamenskiy/kztest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newRandomHC(t *testing.T, bits int) *hashcash.Hashcash {
	resource := make([]byte, 5)
	_, err := rand.Read(resource)
	require.NoError(t, err)
	ret, err := hashcash.New(bits, fmt.Sprintf("%x", resource))
	require.NoError(t, err)
	return ret
}

func TestPOWManager_CheckPOW(t *testing.T) {

	ctx := context.Background()
	lcPOWManager := New(3)

	newhc := func(bits int) *hashcash.Hashcash { return newRandomHC(t, bits) }

	newhcomputed := func(bits int) *hashcash.Hashcash {
		hc := newRandomHC(t, bits)
		err := hc.Compute(0)
		require.NoError(t, err)
		return hc
	}

	kztest.RunTests(t, lcPOWManager.CheckPOW, []kztest.TestKit{
		{Arg1: ctx, Arg2: newhc(0).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhc(1).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhc(2).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhc(3).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhc(4).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhc(5).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(0).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(1).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(2).String(), Result1: false, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(3).String(), Result1: true, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(4).String(), Result1: true, Result2: assert.NoError},
		{Arg1: ctx, Arg2: newhcomputed(5).String(), Result1: true, Result2: assert.NoError},
	})
}
