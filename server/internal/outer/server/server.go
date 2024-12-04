package server

import (
	"context"
	"errors"
	"net"

	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
)

var ErrCloseConnection = errors.New("close connection")

type Ctx = context.Context

const MessageSizeBytesLength = 4

type Handler interface {
	HandleRequest(Ctx, string, *proto.Request) (*proto.Response, error)
}

type Server struct {
	lgr     *lg.Logger
	Addr    string
	Handler Handler

	listener net.Listener

	wrkCtx        context.Context
	wrkCtxCancelF context.CancelFunc
}

func New(addr string, handler Handler) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
	}
}

func (s *Server) ListenAndHandle(ctx Ctx) (err error) {
	s.wrkCtx, s.wrkCtxCancelF = context.WithCancel(ctx)

	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return lge.WrapErrWithCaller(err)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.lgr.Error(ctx, "listener.Accept() error", lga.Err(err))

			sWrkCtxErr := s.wrkCtx.Err()
			if sWrkCtxErr != nil { // it means context is finished
				s.lgr.Error(ctx, "sWrkCtxErr non nil", lga.Any("sWrkCtxErr", sWrkCtxErr))
				return nil
			}

			continue
		}
		go s.HandleConnection(s.wrkCtx, conn)
	}
}

func (s *Server) HandleConnection(ctx Ctx, conn net.Conn) {
	defer func() {
		connCloseErr := conn.Close()
		if connCloseErr != nil {
			s.lgr.Error(ctx, "conn.Close() error", lga.Any("connCloseErr", connCloseErr))
		}
	}()
	var (
		err  error
		req  *proto.Request
		resp *proto.Response
	)
	for {
		req = new(proto.Request)
		err = proto.ReadMessage(ctx, conn, req)
		if err != nil {
			s.lgr.Error(ctx, lge.WrapWithCaller(err), lga.Any("req", req), lga.Err(err))
			req = nil
		}

		lg.Debug(ctx, "req", lga.Any("req", req))

		if req == nil {
			s.lgr.Error(ctx, "nil request leads to error response", lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
			resp = NewUnexpectedServerErrorResponse()
		} else {
			resp, err = s.Handler.HandleRequest(ctx, conn.RemoteAddr().String(), req)
			if err != nil {
				if errors.Is(err, ErrCloseConnection) {
					// defer conn.Close() close the connection
					return
				}
				s.lgr.Error(ctx, "s.HandleRequest error", lga.Any("req", resp), lga.Any("resp", resp), lga.Err(err))
				resp = NewUnexpectedServerErrorResponse()
			}
			if resp == nil {
				s.lgr.Error(ctx, "s.HandleRequest return nil response", lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
				resp = NewUnexpectedServerErrorResponse()
			}
		}

		lg.Debug(ctx, "resp", lga.Any("resp", resp))

		for retry := 0; retry < 3; retry++ {
			err = proto.SendMessage(ctx, conn, resp)
			if err != nil {
				s.lgr.Error(ctx, "proto.SendMessage err", lga.Any("retry", retry), lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
				continue
			}
			break
		}
	}
}

func (s *Server) Close(ctx Ctx) {
	s.wrkCtxCancelF()

	err := s.listener.Close()
	if err != nil {
		s.lgr.Error(ctx, "listener.Close() error", lga.Err(err))
	}

	return
}

func NewUnexpectedServerErrorResponse() (ret *proto.Response) {
	ret = new(proto.Response)
	ret.Type = proto.Response_ERROR
	ret.Resp = &proto.Response_Error{
		Error: &proto.Error{
			Code:    proto.Error_UNEXPECTED_INTERNAL_ERROR,
			Message: "unexpected internal server error occured",
		},
	}
	return ret
}
