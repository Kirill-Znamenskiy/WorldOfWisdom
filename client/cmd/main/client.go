package main

import (
	"context"
	"fmt"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/hashcash"
	"net"
	"os"
	"time"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/client/internal/config"
)

type Ctx = context.Context

//nolint:gochecknoglobals // build app version
var prvBuildGitShowVersion = "UNKNOWN"

func main() {
	ctx := context.Background()

	lg.IsTryExtractWrkLoggerFromCtx = false
	lg.DefaultLogger = lg.MustNewLogger(lg.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Init(ctx)
	if err != nil {
		lg.Error(ctx, lge.WrapWithCaller(err))
		os.Exit(1)
	}
	cfg.BuildGitShowVersion = prvBuildGitShowVersion
	lg.Info(ctx, "Config successfully inited.", lga.String("cfg.BuildGitShowVersion", cfg.BuildGitShowVersion))

	err = lg.Wrk(ctx).ParseAndSetLevel(cfg.LogLevel)
	lg.Info(ctx,
		"lg.Default().ParseAndSetLevel(cfg.LogLevel)",
		lga.String("cfg.LogLevel", cfg.LogLevel),
		lga.Any("err-result", err),
	)
	if err != nil {
		lg.Error(ctx, lge.WrapWithCaller(err, lga.Str("cfg.LogLevel", cfg.LogLevel)))
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	conn, err := net.Dial("tcp", cfg.ServerAddress)
	if err != nil {
		lg.Error(ctx, "net.Dial", lga.Err(err))
		os.Exit(1)
	}

	for {
		time.Sleep(5 * time.Second)

		err = run(ctx, cfg, conn)
		if err != nil {
			lg.Error(ctx, "run", lga.Err(err))
		}
		// defer conn.Close()
	}
}

func run(ctx Ctx, cfg *config.Config, conn net.Conn) (err error) {
	req := new(proto.Request)
	req.Type = proto.Request_WISDOM_REQUEST

	err = proto.SendMessage(ctx, conn, req)
	if err != nil {
		return err
	}

	resp := new(proto.Response)
	err = proto.ReadMessage(ctx, conn, resp)
	if err != nil {
		return err
	}

	if resp.Type == proto.Response_ERROR {
		if resp.GetError().GetCode() == proto.Error_INVALID_POW {
			hc, err := hashcash.Parse(resp.GetChallenge())
			if err != nil {
				return err
			}
			err = hc.Compute(cfg.POWMaxAttempts)
			if err != nil {
				return err
			}

			req := new(proto.Request)
			req.Type = proto.Request_WISDOM_REQUEST
		}
	}

	switch resp.Type {
	case proto.Response_ERROR:
	case proto.Response_WISDOM_RESPONSE:
		fmt.Printf("WISDOM QUOTE: %s", resp.GetWisdomResponse().Quote)
	default:
		return lge.New("unknown response type", lga.Any("resp", resp))
	}

	return nil
}
