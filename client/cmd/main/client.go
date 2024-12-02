package main

import (
	"context"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
	"os"

	"github.com/Kirill-Znamenskiy/kzlogger/lg"
)

type Ctx = context.Context

//nolint:gochecknoglobals // build app version
var prvBuildGitShowVersion = "UNKNOWN"

func main() {
	ctx := context.Background()

	lg.IsTryExtractWrkLoggerFromCtx = false
	lg.DefaultLogger = lg.MustNewLogger(lg.NewTextHandler(os.Stdout, nil))

	req := proto.Request{}

	_ = prvBuildGitShowVersion
}
