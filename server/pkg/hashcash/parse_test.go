package hashcash

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {

	hc, err := New(11, "any-resource:1351:asdgsa:127.0.0.1")
	require.NoError(t, err)

	parsed, err := Parse(hc.String())
	require.NoError(t, err)

	require.Equal(t, hc, parsed)

}
