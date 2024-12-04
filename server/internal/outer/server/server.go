package server

import (
	"context"
	"errors"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"
	"github.com/sourcegraph/conc/panics"
	"io"
	"net"
	"time"
)

var ErrCloseConnection = errors.New("close connection")

type Ctx = context.Context

const (
	MessageSizeBytesLength  = 4
	SendMessageRetriesCount = 3
)

type Handler interface {
	HandleRequest(Ctx, string, *proto.Request) (*proto.Response, error)
}

type Server struct {
	lgr     lg.LoggerInterface
	Addr    string
	Handler Handler

	wrkCtx        context.Context
	wrkCtxCancelF context.CancelFunc
}

func New(logger lg.LoggerInterface, addr string, handler Handler) *Server {
	return &Server{
		lgr:     logger,
		Addr:    addr,
		Handler: handler,
	}
}

func (s *Server) ListenAndHandle(ctx Ctx) (err error) {
	s.wrkCtx, s.wrkCtxCancelF = context.WithCancel(ctx)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return lge.WrapErrWithCaller(err)
	}

	go func() {
		select {
		case <-ctx.Done():
			err := listener.Close()
			if err != nil {
				s.lgr.Error(ctx, "listener.Close() error", lga.Err(err))
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.lgr.Error(ctx, "listener.Accept() error", lga.Err(err))
			sWrkCtxErr := s.wrkCtx.Err()
			if sWrkCtxErr != nil { // it means context is finished
				s.lgr.Error(ctx, "sWrkCtxErr non nil", lga.Any("sWrkCtxErr", sWrkCtxErr))
				return nil
			}

			continue
		}
		s.lgr.Info(ctx, "start handle connection", lga.String("conn", conn.LocalAddr().String()+"<=>"+conn.RemoteAddr().String()))
		go s.HandleConnection(s.wrkCtx, conn)
	}
}

func (s *Server) HandleConnection(ctx Ctx, conn net.Conn) {
	go func() {
		select {
		case <-ctx.Done():
			s.closeConnection(ctx, conn)
		}
	}()
	defer s.closeConnection(ctx, conn)
	var (
		err  error
		req  *proto.Request
		resp *proto.Response
	)
	for {
		req = new(proto.Request)
		err = proto.ReadMessage(conn, req)
		if err != nil {
			s.lgr.Error(ctx, "proto.ReadMessage error", lga.Err(err), lga.Any("req", req))
			if errors.Is(err, io.EOF) || s.wrkCtx.Err() != nil {
				return
			}
			req = nil
		}

		lg.Debug(ctx, "req", lga.Any("req", req))

		if req == nil {
			s.lgr.Error(ctx, "nil request leads to error response", lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
			resp = NewUnexpectedServerErrorResponse()
		} else {
			recoveredPanic := panics.Try(func() {
				resp, err = s.Handler.HandleRequest(ctx, conn.RemoteAddr().String(), req)
			})
			if recoveredPanic != nil {
				s.lgr.Error(ctx, " s.Handler.HandleRequest panic", lga.Any("panic", recoveredPanic.AsError()), lga.Any("resp", resp), lga.Err(err))
				err = nil
				resp = NewUnexpectedServerErrorResponse()
			}
			if err != nil {
				s.lgr.Error(ctx, " s.Handler.HandleRequest error", lga.Any("req", resp), lga.Any("resp", resp), lga.Err(err))
				if errors.Is(err, ErrCloseConnection) {
					return
				}
				resp = NewUnexpectedServerErrorResponse()
			}
			if resp == nil {
				s.lgr.Error(ctx, "s.HandleRequest return nil response", lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
				resp = NewUnexpectedServerErrorResponse()
			}
		}

		lg.Debug(ctx, "resp", lga.Any("resp", resp))

		for retry := 0; retry < SendMessageRetriesCount; retry++ {
			err = proto.SendMessage(conn, resp)
			if err != nil {
				s.lgr.Error(ctx, "proto.SendMessage err", lga.Any("retry", retry), lga.Any("req", req), lga.Any("resp", resp), lga.Err(err))
				if errors.Is(err, io.EOF) || s.wrkCtx.Err() != nil {
					return
				}
				continue
			}
			break
		}
	}
}

func (s *Server) Close(ctx Ctx) {
	s.wrkCtxCancelF()
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

func (s *Server) closeConnection(ctx Ctx, conn net.Conn) {
	s.lgr.Info(ctx, "close connection", lga.String("conn", conn.LocalAddr().String()+"<=>"+conn.RemoteAddr().String()))
	err := conn.SetDeadline(time.Now())
	if err != nil {
		s.lgr.Error(ctx, "conn.SetDeadline(time.Now()) error", lga.Err(err))
	}
	err = conn.Close()
	if err != nil {
		s.lgr.Error(ctx, "conn.Close() error", lga.Err(err))
	}
}
