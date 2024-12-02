package repo

import (
	"context"
	"errors"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/entity"
)

type Ctx = context.Context

type Repo struct {
	wisdoms []*entity.Wisdom
}

func New(wisdomQuotes []string) (ret *Repo) {
	ret = &Repo{}
	ret.wisdoms = make([]*entity.Wisdom, 0, len(wisdomQuotes))
	for nn, wq := range wisdomQuotes {
		ret.wisdoms = append(ret.wisdoms, &entity.Wisdom{
			NN:    entity.WisdomNN(nn),
			Quote: entity.WisdomQuote(wq),
		})
	}
	return ret
}

func (r *Repo) GetWisdom(ctx Ctx, nn entity.WisdomNN) (ret *entity.Wisdom, err error) {
	if len(r.wisdoms) == 0 {
		return nil, errors.New("no any wisdoms")
	}
	ind := uint32(nn)
	ind = ind % uint32(len(r.wisdoms))

	return r.wisdoms[ind], nil
}
