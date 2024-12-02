package server

import (
	"context"
	"encoding/binary"
	"github.com/Kirill-Znamenskiy/WorldOfWisdom/server/pkg/proto"
	"github.com/Kirill-Znamenskiy/kzlogger/lg"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"
	protobuf "google.golang.org/protobuf/proto"
	"io"
	"net"
	"unsafe"
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
	var (
		bs   []byte
		size uint32
	)
	for {
		bs = make([]byte, unsafe.Sizeof(size))
		_, err = io.ReadFull(conn, bs)
		if err != nil {
			return lge.WrapWithCaller(err)
		}
		size = binary.BigEndian.Uint32(bs)

		bs = make([]byte, size)
		_, err = io.ReadFull(conn, bs)
		if err != nil {
			return lge.WrapWithCaller(err)
		}

		req := new(proto.Request)
		err = protobuf.Unmarshal(bs, req)
		if err != nil {
			//return err
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

		err = SendResponse(ctx, conn, resp)
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

func SendResponse(ctx Ctx, conn net.Conn, resp *proto.Response) error {
	respBs, err := protobuf.Marshal(resp)
	if err != nil {
		return err
	}

	size := uint32(len(respBs))
	sizeBs := binary.BigEndian.AppendUint32(nil, size)

	bs := make([]byte, 0, len(sizeBs)+len(respBs))
	bs = append(bs, sizeBs...)
	bs = append(bs, respBs...)

	n, err := conn.Write(bs)
	if err != nil {
		return err
	}
	if n != len(bs) {
		return io.ErrShortWrite
	}

	return nil
}
