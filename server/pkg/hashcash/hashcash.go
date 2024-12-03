package hashcash

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
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
	ZeroBitsMaxValue       = 32
	MaxAttemtptsMaxValue   = 1<<32 - 1
	RandBytesDefaultLength = 16

	ZeroByte byte = '0'
)

var (
	ErrEmptyResource     = errors.New("empty resource")
	ErrResourceIsTooLong = errors.New("resource is too long")

	ErrZeroBitsIsTooBig   = errors.New("zero bits is too big")
	ErrZeroBitsIsNegative = errors.New("zero bits is negative")

	ErrMaxAttemtpsIsTooBig = errors.New("max attempts is too big")

	ErrComputingMaxAttemptsExceeded = errors.New("computing max attempts exceeded")
)

// Hashcash
// see https://en.wikipedia.org/wiki/Hashcash
type Hashcash struct {
	ver       uint8
	bits      uint8
	date      time.Time
	resource  string
	extension string
	rand      []byte
	counter   uint64
}

func New(bits int, resource string) (ret *Hashcash, err error) {
	if bits < 0 {
		return nil, ErrZeroBitsIsNegative
	}
	if bits > ZeroBitsMaxValue {
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
		bits:      uint8(bits),
		date:      time.Now().UTC().Truncate(time.Second),
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

func (hc *Hashcash) GetBits() uint8 {
	return hc.bits
}
func (hc *Hashcash) GetDate() time.Time {
	return hc.date
}

func (hc *Hashcash) String() string {
	cbs := []byte(strconv.FormatUint(hc.counter, 10))
	return fmt.Sprintf("%d:%d:%d:%s:%s:%s:%s",
		hc.ver,
		hc.bits,
		hc.date.Unix(),
		base64.StdEncoding.EncodeToString([]byte(hc.resource)),
		hc.extension,
		base64.StdEncoding.EncodeToString(hc.rand),
		base64.StdEncoding.EncodeToString(cbs),
	)
}

func (hc *Hashcash) IsCorrect() bool {
	return IsCorrect(hc.String(), int(hc.bits))
}

func (hc *Hashcash) Compute(maxAttempts uint64) error {
	if hc.IsCorrect() {
		return nil
	}
	if maxAttempts > MaxAttemtptsMaxValue {
		return ErrMaxAttemtpsIsTooBig
	}
	hc.counter = 0
	for maxAttempts == 0 || hc.counter <= maxAttempts {
		if hc.IsCorrect() {
			return nil
		}
		hc.counter++
	}
	return ErrComputingMaxAttemptsExceeded
}

func IsCorrect(hc string, bits int) bool {
	if hc == "" {
		return false
	}
	if bits == 0 {
		return true
	}
	if bits < 0 {
		return false
	}

	hash := CalcHashSum(hc)
	if bits > len(hash) {
		return false
	}

	for _, b := range hash[:bits] {
		if b != ZeroByte {
			return false
		}
	}

	return true
}

func CalcHashSum[T string | []byte](data T) (ret []byte) {
	hash := sha1.Sum([]byte(data))
	ret = make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(ret, hash[:])
	return ret
}
