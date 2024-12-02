package hashcash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"
)

const (
	Version                = 1
	ResourceMaxLength      = 1024
	ZeroBitsMax            = 32
	RandBytesDefaultLength = 16

	ZeroByte byte = '0'
)

var (
	ErrEmptyResource     = errors.New("empty resource")
	ErrResourceIsTooLong = errors.New("resource is too long")
	ErrEmptyZeroBits     = errors.New("empty zero bits")
	ErrZeroBitsIsTooBig  = errors.New("zero bits is too big")
)

// Hashcash
// see https://en.wikipedia.org/wiki/Hashcash
type Hashcash struct {
	ver       uint
	bits      uint
	date      time.Time
	resource  string
	extension string
	rand      []byte
	counter   uint64
}

func New(bits uint, resource string) (ret *Hashcash, err error) {
	if bits <= 0 {
		return nil, ErrEmptyZeroBits
	}
	if bits > ZeroBitsMax {
		return nil, ErrZeroBitsIsTooBig
	}
	if resource == "" {
		return nil, ErrEmptyResource
	}
	if len(resource) > ResourceMaxLength {
		return nil, lge.Wrap(ErrResourceIsTooLong, lga.Int("len(resource)", len(resource)))
	}

	ret = &Hashcash{
		ver:       Version,
		bits:      bits,
		date:      time.Now().UTC(),
		resource:  resource,
		extension: "",
		rand:      make([]byte, RandBytesDefaultLength),
		counter:   0,
	}

	_, err = rand.Read(ret.rand)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (h *Hashcash) String() string {
	cbs := []byte(strconv.FormatUint(h.counter, 10))
	return fmt.Sprintf("%d:%d:%d:%s:%s:%s:%s",
		h.ver,
		h.bits,
		h.date.Unix(),
		h.resource,
		h.extension,
		base64.StdEncoding.EncodeToString(h.rand),
		base64.StdEncoding.EncodeToString(cbs),
	)
}

//func (h *Hashcash) Compute(maxAttempts int) error {
//	if maxAttempts > 0 {
//		h.counter = 0
//		for h.counter <= maxAttempts {
//			ok, err := h.Header().IsHashCorrect(h.Bits)
//			if err != nil {
//				return err
//			}
//			if ok {
//				return nil
//			}
//			h.counter++
//		}
//	}
//
//	return ErrComputingMaxAttemptsExceeded
//}
//
//// Key - returns string presentation of hashcash without counter
//// Key is using to match original hashcash with solved hashcash
//func (h *Hashcash) Key() string {
//	return fmt.Sprintf("%d:%d:%s:%d", h.Bits, h.Date.Unix(), h.Resource, binary.BigEndian.Uint32(h.Rand))
//}
//
//// Header - returns string presentation of hashcash to share it
//func (h *Hashcash) Header() Header {
//	return Header(fmt.Sprintf("1:%d:%s:%s:%s:%s:%s",
//		h.Bits,
//		h.Date.Format(time.r),
//		h.Resource,
//		h.Extension,
//		base64.StdEncoding.EncodeToString(h.Rand),
//		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
//	))
//}
//
//// ParseHeader - parse hashcah from header
//func ParseHeader(header string) (hashcash *Hashcash, err error) {
//	parts := strings.Split(header, ":")
//
//	if len(parts) < 7 {
//		return nil, ErrIncorrectHeaderFormat
//	}
//	if len(parts) > 7 {
//		for i := 0; i < len(parts)-7; i++ {
//			parts[3] += ":" + parts[3+i+1]
//		}
//		parts[4] = parts[len(parts)-3]
//		parts[5] = parts[len(parts)-2]
//		parts[6] = parts[len(parts)-1]
//		parts = parts[:7]
//	}
//	if parts[0] != "1" {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	hashcash = &Hashcash{}
//
//	hashcash.Bits, err = strconv.Atoi(parts[1])
//	if err != nil {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	hashcash.Date, err = time.ParseInLocation(dateLayout, parts[2], time.UTC)
//	if err != nil {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	hashcash.Resource = parts[3]
//	hashcash.Extension = parts[4]
//
//	hashcash.Rand, err = base64.StdEncoding.DecodeString(parts[5])
//	if err != nil {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	counterStr, err := base64.StdEncoding.DecodeString(parts[6])
//	if err != nil {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	hashcash.counter, err = strconv.Atoi(string(counterStr))
//	if err != nil {
//		return nil, ErrIncorrectHeaderFormat
//	}
//
//	return
//}
//
//// Header - string presentation of hashcash
//// Format - 1:Bits:Date:Resource:externsion:Rand:counter
//type Header string
//
//// IsHashCorrect - does header hash constain zero Bits enough
//func (header Header) IsHashCorrect(bits int) (ok bool, err error) {
//	if bits <= 0 {
//		return false, ErrZeroBitsMustBeMoreThanZero
//	}
//
//	hash, err := header.sha1()
//	if err != nil {
//		return ok, err
//	}
//	if len(hash) < bits {
//		return false, ErrHashLengthLessThanZeroBits
//	}
//
//	ok = true
//	for _, s := range hash[:bits] {
//		if s != zeroBit {
//			ok = false
//			break
//		}
//	}
//	return
//}
//
//func (header Header) sha1() (hash string, err error) {
//	hasher := sha1.New()
//	if _, err = hasher.Write([]byte(header)); err != nil {
//		return
//	}
//
//	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
//}
