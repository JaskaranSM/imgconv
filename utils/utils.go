package utils

import (
	"bytes"
	"image"
	"image/png"

	"github.com/strukturag/libheif/go/heif"
)

func GetLibHeifVersion() string {
	return heif.GetVersion()
}

func DecodeHEIFImageBytes(inp []byte) (image.Image, error) {
	c, err := heif.NewContext()
	if err != nil {
		return nil, err
	}
	err = c.ReadFromMemory(inp)
	if err != nil {
		return nil, err
	}
	handle, err := c.GetPrimaryImageHandle()
	if err != nil {
		return nil, err
	}
	imgHeif, err := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if err != nil {
		return nil, err
	}
	return imgHeif.GetImage()
}

func ConvertImageToPngBytes(inp []byte) ([]byte, error) {
	img, err := DecodeHEIFImageBytes(inp)
	var out bytes.Buffer
	err = png.Encode(&out, img)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
