package mWisdom

import (
	"context"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/entity"
)

type Ctx = context.Context

type Repo interface {
	GetWisdom(Ctx, entity.WisdomNN) (*entity.Wisdom, error)
}

type Randomizer interface {
	Uint32() uint32
}

type WisdomManager struct {
	repo       Repo
	randomizer Randomizer
}

func New(repo Repo, randomizer Randomizer) *WisdomManager {
	return &WisdomManager{
		repo:       repo,
		randomizer: randomizer,
	}
}

func (m *WisdomManager) GetRandomWisdomQuote(ctx Ctx) (ret entity.WisdomQuote, err error) {
	randNN := entity.WisdomNN(m.randomizer.Uint32())
	wisdom, err := m.repo.GetWisdom(ctx, randNN)
	if err != nil {
		return "", err
	}
	return wisdom.Quote, nil
}
