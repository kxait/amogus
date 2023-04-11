package config

import (
	"math"
	"os"
)

const DefaultPartLength int64 = 32000

type HashesInfo struct {
	Parts      int64
	PartLength int64
}

func GetHashesInfo(filename string) (*HashesInfo, error) {
	st, err := os.Stat(filename)

	if err != nil {
		return nil, err
	}

	return &HashesInfo{
		Parts:      int64(math.Ceil(float64(st.Size()) / float64(DefaultPartLength))),
		PartLength: DefaultPartLength,
	}, nil
}

func GetHashesPart(filename string, partNo int64) ([]byte, error) {
	st, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	offset := DefaultPartLength * partNo
	lengthToEnd := st.Size() - offset
	var bufsize int
	if lengthToEnd < DefaultPartLength {
		bufsize = int(lengthToEnd)
	} else {
		bufsize = int(DefaultPartLength)
	}

	file.Seek(DefaultPartLength*partNo, 0)

	buf := make([]byte, bufsize)

	file.Read(buf)

	return buf, nil
}
