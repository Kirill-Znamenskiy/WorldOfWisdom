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

func TestPOWManager_CheckPOW(t *testing.T) {

	ctx := context.Background()
	lcPOWManager := New(3)

	for i := 0; i < 3; i++ {
		fmt.Printf("\n\n\n")
		hc, _ := hashcash.New(1, "asdfa")
		fmt.Printf("before-hc: %s\n", hc.String())
		hc.Compute(1)
		fmt.Printf("after-hc: %s\n", hc.String())
	}

	hc := func(bits int) *hashcash.Hashcash {
		resource := make([]byte, 5)
		_, err := rand.Read(resource)
		require.NoError(t, err)
		ret, err := hashcash.New(bits, string(resource))
		require.NoError(t, err)
		return ret
	}

	kztest.RunTests(t, lcPOWManager.CheckPOW, []kztest.TestKit{
		{Arg1: ctx, Arg2: hc(3).String(), Result1: false, Result2: assert.NoError},
	})
}
