package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/config"
	mPOW "github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/inner/managers/pow-manager"
	mWisdom "github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/inner/managers/wisdom-manager"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/outer/handlers"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/outer/repo"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/outer/server"
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

	lcPOWManager := mPOW.New(cfg.Server.POW.ZeroBitsCount)

	lcRepo := repo.New(prvWisdomQuotes)

	lcRandomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	lcWisdomManager := mWisdom.New(lcRepo, lcRandomizer)

	lcHandlers := handlers.New(lcPOWManager, lcWisdomManager)

	lcServerInst := server.New(cfg.Server.Address, lcHandlers)

	finishedCh := make(chan error, 3)

	ctx, ctxCancelF := context.WithCancel(ctx)

	go func() {
		lg.Info(ctx, "TCP Server: go Listen and Handle")
		finishedCh <- lcServerInst.ListenAndHandle(ctx)
	}()

	closeFunc := func(timeout time.Duration) {
		ctxCancelF()

		lg.Info(ctx, "TCP Server: now graceful stop")
		bgCtx, bgCtxCancelF := context.WithTimeout(context.Background(), timeout)
		defer bgCtxCancelF()
		err = lcServerInst.Close(bgCtx)
		if err != nil {
			lg.Error(bgCtx, lge.WrapWithCaller(err))
		}
	}

	forKillSignalsCh := make(chan os.Signal, 1)
	signal.Notify(forKillSignalsCh, os.Kill)
	forInterruptSignalsCh := make(chan os.Signal, 1)
	signal.Notify(forInterruptSignalsCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sg := <-forKillSignalsCh:
		lg.Error(ctx, "Kill signal accepted!", lga.Any("sg", sg))
		closeFunc(0)
	case sg := <-forInterruptSignalsCh:
		lg.Error(ctx, "Interrupt signal accepted!", lga.Any("sg", sg))
		closeFunc(3 * time.Second)
	case err = <-finishedCh:
		lg.Error(ctx, "Somebody has finished with error!", lga.Err(err))
		closeFunc(3 * time.Second)
	}
}
