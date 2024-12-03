package hashcash

import (
	"encoding/base64"
	"errors"
	"github.com/Kirill-Znamenskiy/kzlogger/lga"
	"github.com/Kirill-Znamenskiy/kzlogger/lge"
	"strconv"
	"strings"
	"time"
)

var ErrIncorrectFormat = errors.New("incorrect hashcash format")

func Parse(hc string) (ret *Hashcash, err error) {
	parts := strings.Split(hc, ":")

	if len(parts) != 7 {
		return nil, ErrIncorrectFormat
	}

	ret = new(Hashcash)

	if parts[0] != strconv.Itoa(Version) {
		return nil, ErrIncorrectFormat
	}
	ret.ver = Version

	bits, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, ErrIncorrectFormat
	}
	if bits < 0 {
		return nil, ErrZeroBitsIsNegative
	}
	if bits > ZeroBitsMaxValue {
		return nil, ErrZeroBitsIsTooBig
	}
	ret.bits = uint8(bits)

	unixsecs, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || unixsecs <= 0 {
		return nil, ErrIncorrectFormat
	}
	ret.date = time.Unix(unixsecs, 0).UTC()

	resourceBs, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, ErrIncorrectFormat
	}
	resource := string(resourceBs)
	if resource == "" {
		return nil, ErrEmptyResource
	}
	if len(resource) > ResourceMaxLength {
		return nil, lge.Wrap(ErrResourceIsTooLong, lga.Int("len(resource)", len(resource)))
	}
	ret.resource = resource

	ext := parts[4]
	if ext != "" {
		return nil, ErrIncorrectFormat
	}
	ret.extension = ext

	rand, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, ErrIncorrectFormat
	}
	ret.rand = rand

	counterBs, err := base64.StdEncoding.DecodeString(parts[6])
	if err != nil {
		return nil, ErrIncorrectFormat
	}
	counter, err := strconv.ParseUint(string(counterBs), 10, 64)
	if err != nil {
		return nil, ErrIncorrectFormat
	}
	ret.counter = counter

	return ret, nil
}
