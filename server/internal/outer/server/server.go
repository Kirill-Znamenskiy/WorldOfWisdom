package server

import (
	"context"
	"net"

	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"

	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
)

type Ctx = context.Context

const MessageSizeBytesLength = 4

type Handler interface {
	HandleRequest(Ctx, string, *proto.Request) (*proto.Response, error)
}

type Server struct {
	lgr     *lg.Logger
	Addr    string
	Handler Handler
}

func New(addr string, handler Handler) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
	}
}

func (s *Server) ListenAndHandle(ctx Ctx) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return lge.WrapErrWithCaller(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return lge.WrapErrWithCaller(err)
		}
		go s.GoHandleConnection(ctx, conn)
	}
	return nil
}

func (s *Server) GoHandleConnection(ctx Ctx, conn net.Conn) {
	err := s.HandleConnection(ctx, conn)
	if err != nil {
		s.lgr.Error(ctx, "s.HandleConnection(ctx, conn)", lga.Err(err))
	}
}

func (s *Server) HandleConnection(ctx Ctx, conn net.Conn) (err error) {
	defer conn.Close()
	for {
		req := new(proto.Request)
		err = proto.ReadMessage(ctx, conn, req)
		if err != nil {
			return err
		}
		lg.Debug(ctx, "req", lga.Any("req", req))

		resp, err := s.HandleRequest(ctx, conn.RemoteAddr().String(), req)
		if err != nil {
			return err
		}
		if resp == nil {
			resp = new(proto.Response)
		}

		lg.Debug(ctx, "resp", lga.Any("resp", resp))

		err = proto.SendMessage(ctx, conn, resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) HandleRequest(ctx Ctx, client string, req *proto.Request) (resp *proto.Response, err error) {
	return s.Handler.HandleRequest(ctx, client, req)
}

func (s *Server) Close(ctx Ctx) (err error) {
	return nil
}
