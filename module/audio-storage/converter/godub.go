package converter

import (
	"bytes"

	gc "github.com/iFaceless/godub/converter"
)

type GoDub struct {
}

func New() *GoDub {
	return &GoDub{}
}

func (gd *GoDub) Convert(source []byte, extension string) ([]byte, error) {
	writer := new(bytes.Buffer)
	err := gc.NewConverter(writer).
		WithBitRate(64000).
		WithDstFormat(extension).
		Convert(bytes.NewReader(source))

	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
