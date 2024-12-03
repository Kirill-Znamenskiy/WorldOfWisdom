package handlers

import (
	"context"
	"errors"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/entity"
	mPOW "github.com/Kirill-Znamenskiy/WorldOfWisdom/server/internal/inner/managers/pow-manager"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
)

type Ctx = context.Context

type POWManager interface {
	CheckPOW(Ctx, mPOW.POW) (bool, error)
	GenerateNewChallenge(Ctx, mPOW.Client) (mPOW.Challenge, error)
}

type WisdomManager interface {
	GetRandomWisdomQuote(Ctx) (entity.WisdomQuote, error)
}

type Handlers struct {
	prvPOWManager    POWManager
	prvWisdomManager WisdomManager
}

func New(pPOWManager POWManager, pWisdomManager WisdomManager) *Handlers {
	return &Handlers{
		prvPOWManager:    pPOWManager,
		prvWisdomManager: pWisdomManager,
	}
}

func (hs *Handlers) HandleRequest(ctx Ctx, client string, req *proto.Request) (resp *proto.Response, err error) {
	if req.Type == proto.RequestType_REQ_QUIT {
		return nil, errors.New("close due to quit request")
	}
	resp = &proto.Response{}
	resp.Challenge, err = hs.prvPOWManager.GenerateNewChallenge(ctx, client)
	if err != nil {
		return nil, err
	}

	var ok bool
	switch req.Type {
	case proto.RequestType_WISDOM_REQUEST:
		ok, err = hs.prvPOWManager.CheckPOW(ctx, req.Pow)
		if err != nil {
			return nil, err
		}
		if !ok {
			resp.Type = proto.ResponseType_ERROR
			resp.Resp = &proto.Response_Error{
				Error: &proto.Error{
					Code:    101,
					Message: "provide valid proof of work",
				},
			}
		}
		resp.Type = proto.ResponseType_WISDOM_RESPONSE
		lcGetWisdomResponse, err := hs.HandleWisdomRequest(ctx, req.GetWisdomRequest())
		if err != nil {
			return nil, err
		}
		resp.Resp = &proto.Response_WisdomResponse{
			WisdomResponse: lcGetWisdomResponse,
		}
	default:
		resp.Type = proto.ResponseType_ERROR
		resp.Resp = &proto.Response_Error{
			Error: &proto.Error{
				Code:    100,
				Message: "unknown request unit type",
			},
		}
	}

	return resp, nil
}
