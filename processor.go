package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"image/jpeg"
	"image/png"

	"net/http"
)

type File struct {
	name   string
	buffer []byte
}

func readFile(name string) ([]byte, error) {
	buf, err := os.ReadFile(name)
	return buf, err
}

func (f *File) convertToJPG() ([]byte, error) {
	imageFormat := http.DetectContentType(f.buffer)

	switch imageFormat {
	case "image/png":
		imageBuffer, err := png.Decode(bytes.NewReader(f.buffer))
		if err != nil {
			return nil, fmt.Errorf("unable to decode, %v", err)
		}

		outBuffer := new(bytes.Buffer)
		if err := jpeg.Encode(outBuffer, imageBuffer, nil); err != nil {
			return nil, fmt.Errorf("unable to encode, %v", err)
		}

		return outBuffer.Bytes(), nil
	case "image/jpeg", "image/jpg":
		return f.buffer, nil
	default:
		return nil, errors.New("unsupported file format")
	}
}

func (f *File) compress() ([]byte, error) {
	image, err := jpeg.Decode(bytes.NewReader(f.buffer))
	if err != nil {
		return nil, fmt.Errorf("unable to decode image, %v", err)
	}

	compressedImage := new(bytes.Buffer)
	err = jpeg.Encode(compressedImage, image, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return nil, fmt.Errorf("unable to compress image, %v", err)
	}

	return compressedImage.Bytes(), nil
}

func uploadToS3(string) error {
	// TODO: Implementations
	return nil
}
