package util

import (
	"encoding/json"
	"strings"
)

// ImageData ...
type ImageData struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// SaveImageData save a data uri to a file and return a url where the file can be reached
func SaveImageData(env *Environment, imgData json.RawMessage) (image *ImageData, err error) {
	saver := &AssetSave{}
	saver.Init(env)

	asset := &Asset{}
	image = &ImageData{}

	// convert images in json
	if len(imgData) > 0 {
		err = json.Unmarshal(imgData, image)
		if err != nil {
			return
		}

		if !strings.HasPrefix(image.Data, "data:") {
			// not a data uri
			return
		}

		// extract data uri and save to disk
		asset, err = saver.Save(&image.Data)
		if err != nil {
			return
		}

		image.Data = asset.FileURL
	}

	return
}
