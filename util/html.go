package util

//
// html.go
// a collection of html image helper functions
// Copyright 2017 Akinmayowa Akinyemi
//

// cspell: ignore Puerkito, goquery

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/xid"
)

// DataURISaver ...
type DataURISaver interface {
	Save(data *string) (*Asset, error)
}

// Asset ...
type Asset struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OwnerID      string `json:"owner_id"`
	AssetType    int    `json:"asset_type"`
	FileName     string `json:"file_name"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// NewAsset creates an asset record and initializes the struct
func NewAsset(name, id string) *Asset {
	asset := &Asset{
		ID:   xid.New().String(),
		Name: name,
	}

	if len(id) > 0 {
		asset.ID = id
	}

	asset.FileName = asset.ID + filepath.Ext(name)

	return asset
}

// AssetSave implements the DataURISaver interface
type AssetSave struct {
	assetFolder  string
	assetBaseURL string
}

// Init ...
func (s *AssetSave) Init(env *Environment) {
	s.assetFolder = env.Paths["upload"]
	s.assetBaseURL = env.Urls["base_url"] + env.Urls["upload"]
}

// AssetFolder returns the asset folder
func (s AssetSave) AssetFolder() string {
	return s.assetFolder
}

// AssetBaseURL the assets base url
func (s AssetSave) AssetBaseURL() string {
	return s.assetBaseURL
}

// Save convert and save as asset datauri in @data
func (s AssetSave) Save(data *string) (*Asset, error) {

	if strings.HasPrefix(*data, ImgDataURIHeader) == false {
		return nil, nil
	}

	asset := NewAsset("", "")
	asset.Name = asset.ID

	fileExt, err := GetDataURIType([]byte(*data))
	if err != nil {
		return nil, err
	}

	asset.FileName += "." + fileExt
	assetFileName := filepath.Join(s.assetFolder, asset.FileName)
	asset.FileURL = s.assetBaseURL + "/" + asset.FileName
	asset.AssetType = 1

	if err := SaveDataURI([]byte(*data), assetFileName); err != nil {
		return nil, err
	}

	return asset, nil
}

// ConvertImgSrcAttr ...
func ConvertImgSrcAttr(html string, saver DataURISaver) (string, error) {
	var (
		err     error
		tags    *goquery.Document
		fileURL string
		asset   *Asset
	)

	htmlBuff := bytes.NewBufferString(html)
	tags, err = goquery.NewDocumentFromReader(htmlBuff)
	if err != nil {
		return html, err
	}

	tags.Find("img").EachWithBreak(func(idx int, el *goquery.Selection) bool {

		data, _ := el.Attr("src")

		asset, err = saver.Save(&data)
		if err != nil {
			return false
		} else if asset == nil && err == nil {
			// src doesn't contain a data uri, continue
			return false
		}

		fileURL = asset.FileURL
		el.SetAttr("src", fileURL)

		return true
	})

	if err != nil {
		return html, err
	}

	return tags.Find("body").Html()
}
