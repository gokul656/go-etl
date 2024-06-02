package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

var pngBuffer = createDummyImage("png")
var invalidBuffer = []byte{}

func BenchmarkETL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		extractedData := extract(files, 8)
		convertedData := convert(extractedData, 8)
		compressedData := compress(convertedData, 8)
		save(compressedData)
	}
}

// Utility function to create a dummy PNG image
func createDummyImage(format string) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	buffer := new(bytes.Buffer)

	switch format {
	case "jpg", "jpeg":
		jpeg.Encode(buffer, img, nil)
	case "png":
		png.Encode(buffer, img)
	}

	return buffer.Bytes()
}
