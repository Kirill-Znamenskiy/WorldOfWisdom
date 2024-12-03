package handlers

import (
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
)

func (hs *Handlers) HandleWisdomRequest(ctx Ctx, req *proto.WisdomRequest) (resp *proto.WisdomResponse, err error) {
	wq, err := hs.prvWisdomManager.GetRandomWisdomQuote(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.WisdomResponse{
		Quote: string(wq),
	}, nil
}
