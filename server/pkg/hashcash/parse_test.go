package hashcash

import (
	"encoding/base64"
	"github.com/Kirill-Znamenskiy/kztest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	hc, err := New(11, "any-resource:1351:asdgsa:127.0.0.1/asdgas@adsga.asdag")
	require.NoError(t, err)

	parsed, err := Parse(hc.String())
	require.NoError(t, err)

	require.Equal(t, hc, parsed)

	kztest.RunTests(t, Parse, []kztest.TestKit{
		{Arg1: "", Result1: assert.Nil, Result2: ErrIncorrectFormat},
		{Arg1: "asdgsa", Result1: assert.Nil, Result2: ErrIncorrectFormat},
		{Arg1: "2:asdgsa", Result1: assert.Nil, Result2: ErrIncorrectFormat},
		{Arg1: "2:asdgsa:2352", Result1: assert.Nil, Result2: ErrIncorrectFormat},
		{Arg1: "1:-1:1733244442:YW55LXJlc291cmNlOjEzNTE6YXNkZ3NhOjEyNy4wLjAuMS9hc2RnYXNAYWRzZ2EuYXNkYWc=::OU+qD0jZO2Pk5ORQVZrmAA==:MA==", Result1: assert.Nil, Result2: ErrZeroBitsIsNegative},
		{Arg1: "1:0:1733244442:YW55LXJlc291cmNlOjEzNTE6YXNkZ3NhOjEyNy4wLjAuMS9hc2RnYXNAYWRzZ2EuYXNkYWc=::OU+qD0jZO2Pk5ORQVZrmAA==:MA==", Result1: assert.NotNil, Result2: assert.NoError},
		{Arg1: "1:1:1733244442:YW55LXJlc291cmNlOjEzNTE6YXNkZ3NhOjEyNy4wLjAuMS9hc2RnYXNAYWRzZ2EuYXNkYWc=::OU+qD0jZO2Pk5ORQVZrmAA==:MA==", Result1: assert.NotNil, Result2: assert.NoError},
		{Arg1: "1:2143423626231:1733244442:YW55LXJlc291cmNlOjEzNTE6YXNkZ3NhOjEyNy4wLjAuMS9hc2RnYXNAYWRzZ2EuYXNkYWc=::OU+qD0jZO2Pk5ORQVZrmAA==:MA==", Result1: assert.Nil, Result2: ErrZeroBitsIsTooBig},
		{Arg1: "1:1:1733244442:" + base64.StdEncoding.EncodeToString([]byte(strings.Repeat("asdf", 777))) + "::OU+qD0jZO2Pk5ORQVZrmAA==:MA==", Result1: assert.Nil, Result2: assert.Error},
		{Arg1: hc.String(), Result1: hc, Result2: assert.NoError},
	})
}
