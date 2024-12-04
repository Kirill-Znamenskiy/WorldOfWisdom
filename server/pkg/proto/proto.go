package proto

import (
	"encoding/binary"
	"io"
	"unsafe"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/Kirill-Znamenskiy/kzlogger/lge"
)

func ReadMessage(reader io.Reader, target protobuf.Message) (err error) {
	var size uint32
	bs := make([]byte, unsafe.Sizeof(size))
	_, err = io.ReadFull(reader, bs)
	if err != nil {
		return lge.WrapWithCaller(err)
	}
	size = binary.BigEndian.Uint32(bs)

	bs = make([]byte, size)
	_, err = io.ReadFull(reader, bs)
	if err != nil {
		return lge.WrapWithCaller(err)
	}

	err = protobuf.Unmarshal(bs, target)
	if err != nil {
		return err
	}

	return nil
}

func SendMessage(writer io.Writer, message protobuf.Message) error {
	messageBs, err := protobuf.Marshal(message)
	if err != nil {
		return err
	}

	size := uint32(len(messageBs))
	sizeBs := binary.BigEndian.AppendUint32(nil, size)

	bs := make([]byte, 0, len(sizeBs)+len(messageBs))
	bs = append(bs, sizeBs...)
	bs = append(bs, messageBs...)

	n, err := writer.Write(bs)
	if err != nil {
		return err
	}
	if n != len(bs) {
		return io.ErrShortWrite
	}

	return nil
}
