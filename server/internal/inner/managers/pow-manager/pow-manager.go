package mPOW

import (
	"context"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/hashcash"
)

type (
	Ctx       = context.Context
	POW       = string
	Client    = string
	Challenge = string
)

type POWManager struct{}

func New() *POWManager {
	return &POWManager{}
}

func (m *POWManager) GenerateNewChallenge(ctx Ctx, client string) (Challenge, error) {
	hc, err := hashcash.New(3, client)
	if err != nil {
		return "", err
	}
	return hc.String(), nil
}

func (m *POWManager) CheckPOW(Ctx, POW) (bool, error) {
	return true, nil
}
