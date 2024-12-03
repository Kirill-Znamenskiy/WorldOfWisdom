package mPOW

import (
	"context"
	"time"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/hashcash"
)

type (
	Ctx       = context.Context
	POW       = string
	Client    = string
	Challenge = string
)

type POWManager struct {
	prvZeroBitsCount uint8
}

func New(pZeroBitsCount uint8) *POWManager {
	return &POWManager{
		prvZeroBitsCount: pZeroBitsCount,
	}
}

func (m *POWManager) GenerateNewChallenge(ctx Ctx, client string) (Challenge, error) {
	hc, err := hashcash.New(int(m.prvZeroBitsCount), client)
	if err != nil {
		return "", err
	}
	return hc.String(), nil
}

func (m *POWManager) CheckPOW(ctx Ctx, pow POW) (bool, error) {
	if !hashcash.IsCorrect(pow, int(m.prvZeroBitsCount)) {
		return false, nil
	}
	hc, err := hashcash.Parse(pow)
	if err != nil {
		return false, err
	}
	if hc.GetBits() < m.prvZeroBitsCount {
		return false, nil
	}
	if hc.GetDate().Add(time.Hour).Before(time.Now()) {
		return false, nil
	}
	return true, nil
}
