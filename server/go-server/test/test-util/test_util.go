package testutil

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

// Middle arg is a close function, should be called by calling function
func GetFileFromPath(path string) (*os.File, func() error, error) {
	// Grab example image to process, decode image, then encode as jpg
	file, err := os.Open(path)
	if err != nil {
		return nil, file.Close, fmt.Errorf("failed to open file: %w", err)
	}
	return file, file.Close, nil
}

func DecodeAndEncodeFileAsJpg(file *os.File) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode img to jpg: %w", err)
	}
	return buffer.Bytes(), nil
}

func DecodeAndEncodeBytesAsJpg(byteSlice []byte) ([]byte, error) {
	// Convert bytes received to a jpg and write to a file in cwd
	imgReturnDecode, err := jpeg.Decode(bytes.NewReader(byteSlice))
	if err != nil {
		return nil, fmt.Errorf("failed to decode return image: %w", err)
	}
	buf := new(bytes.Buffer) //var opts jpeg.Options // opts.Quality = 80
	err = jpeg.Encode(buf, imgReturnDecode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode return image to jpg: %w", err)
	}
	return buf.Bytes(), nil
}

// first arg is a close function, should be called by calling function
func WriteBytesToJpgFile(byteSlice []byte, path string) (func() error, error) {
	jpegFile, err := os.Create(path)
	if err != nil {
		return jpegFile.Close, fmt.Errorf("failed to create test.jpg: %v", err)
	}
	_, err = jpegFile.Write(byteSlice)
	if err != nil {
		return jpegFile.Close, fmt.Errorf("failed to write jpg file: %v", err)
	}
	return jpegFile.Close, nil
}
