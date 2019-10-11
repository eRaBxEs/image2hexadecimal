//
// image.go
// a collection of image manipulation functions
// Copyright 2017 Akinmayowa Akinyemi
//

package util

// cspell: ignore Lanczos

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/labstack/gommon/log"
)

// ImgDataURIHeader the string sequence which is prefixed to image type data uris
const ImgDataURIHeader = "data:image/"

// SaveDataURI converts a datauri back to its original format
// the format type is deduced from the data uri header
func SaveDataURI(data []byte, fileName string) error {

	// read the header and check if the data supplied is a valid data uri
	if !bytes.HasPrefix(data, []byte(ImgDataURIHeader)) {
		return fmt.Errorf("data is not a valid datauri")
	}

	src := bytes.Replace(data, []byte(ImgDataURIHeader), []byte(""), 1)
	idx := bytes.Index(src, []byte(";"))
	if idx == -1 {
		return fmt.Errorf("data is not a valid datauri: cant find mime type")
	}

	// get mime type
	mimeType := string(src[0:idx])
	if mimeType == "jpeg" {
		mimeType = "jpg"
	} else if strings.HasPrefix(mimeType, "svg") {
		mimeType = "svg"
	}

	// get image data (in base64)
	src = src[idx+8 : len(src)]
	raw := make([]byte, len(src))
	_, err := base64.StdEncoding.Decode(raw, src)
	if err != nil {
		return err
	}
	raw = bytes.TrimRightFunc(raw, func(r rune) bool {
		if r == 0 {
			return true
		}

		return false
	})

	log.Debug("raw: ", raw)

	if err := ioutil.WriteFile(fileName, raw, 0644); err != nil {
		return err
	}

	return nil
}

// GetDataURIType ...
func GetDataURIType(data []byte) (string, error) {
	// read the header and check if the data supplied is a valid data uri
	if !bytes.HasPrefix(data, []byte(ImgDataURIHeader)) {
		return "", fmt.Errorf("data is not a valid datauri")
	}

	src := bytes.Replace(data, []byte(ImgDataURIHeader), []byte(""), 1)
	idx := bytes.Index(src, []byte(";"))
	if idx == -1 {
		return "", fmt.Errorf("data is not a valid datauri: cant find mime type")
	}

	// get mime type
	mimeType := string(src[0:idx])
	if mimeType == "jpeg" {
		mimeType = "jpg"
	} else if strings.HasPrefix(mimeType, "svg") {
		mimeType = "svg"
	}

	return mimeType, nil
}

// ResizeImage ...
func ResizeImage(srcFile, dstFile string, width, height int) (err error) {

	srcImg, err := imaging.Open(srcFile)
	if err != nil {
		return nil
	}

	dstImg := imaging.Thumbnail(srcImg, width, height, imaging.Lanczos)
	err = imaging.Save(dstImg, dstFile)
	if err != nil {
		return
	}

	return nil
}
