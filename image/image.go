////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package image contains methods to manipulate png image files.
package image

import (
	"errors"
	"strings"
	"time"

	. "github.com/abitofhelp/go-helpers/string"
	. "github.com/abitofhelp/go-helpers/time"
	"google.golang.org/genproto/googleapis/type/latlng"
	"path/filepath"
)

// Type Image contains metadata to associate with the image when it is persisted.
type Image struct {
	// The path of the image file.
	path string

	// Field filename is the filename without any path information.
	filename string

	// Field latlng is the latitude and longitude where the image was taken.
	latlng latlng.LatLng

	// Field createdUtc is the date/time when the image was taken, in UTC.
	createdUtc time.Time
}

// Function New is a factory that creates an initialized Image.
// Parameter path is the path and filename to an image file that will be processed.
// Returns an initialized instance or error.
func NewImageFromPath(path string) (*Image, error) {
	pipeline := &Image{}
	if pipeline == nil {
		return nil, errors.New("failed to create an instance of Image")
	}

	// Determine the filename from the path.
	directory, filename := filepath.Split(path)

	err := pipeline.SetPath(directory)
	if err != nil {
		return nil, err
	}

	pipeline.SetFileName(filename)

	err = pipeline.SetCreatedUtc(Zero())
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

// Function New is a factory that creates an initialized Image.
// Parameter path is the path to an image file that will be processed.
// Parameter filename is the name of the image file to be processed.
// Returns an initialized instance or error.
func New(path string, filename string) (*Image, error) {
	pipeline := &Image{}
	if pipeline == nil {
		return nil, errors.New("failed to create an instance of Image")
	}

	err := pipeline.SetPath(path)
	if err != nil {
		return nil, err
	}

	pipeline.SetFileName(filename)

	err = pipeline.SetCreatedUtc(Zero())
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

// Method Path gets the path from the instance.
func (i Image) Path() string {
	return i.path
}

// Method SetPath sets the value of the path in the instance.
// If there is an error, an error is returned, otherwise nil.
func (i *Image) SetPath(path string) error {

	if path == "" {
		return errors.New("the path cannot be empty")
	}
	path = CleanStringForPlatform(path)
	i.path = path

	return nil
}

// Method FileName gets the name of the image file from the instance.
func (i Image) FileName() string {
	return i.path
}

// Method SetFileName sets the value of the name of the image file in the instance.
// If there is an error, an error is returned, otherwise nil.
func (i *Image) SetFileName(filename string) error {

	if filename == "" {
		return errors.New("the file name cannot be empty")
	}
	filename = CleanStringForPlatform(filename)
	i.filename = filename

	return nil
}

// Method LatLng gets the latitude and longitude for where the image was created.
func (i Image) LatLng() latlng.LatLng {
	return i.latlng
}

// Method SetLatLng sets latitude and longitude for where the image was created.
// If there is an error, an error is returned, otherwise nil.
func (i *Image) SetLatLng(latLong latlng.LatLng) error {
	i.latlng = latLong

	return nil
}

// Method Created gets the UTC date/time when the instance was created.
func (i Image) CreatedUtc() time.Time {
	return i.createdUtc
}

// Method SetCreatedUtc sets the value of the createdUtc in the instance.
// If there is an error, an error is returned, otherwise nil.
func (i *Image) SetCreatedUtc(createdUtc time.Time) error {
	if createdUtc.IsZero() {
		return errors.New("the createdUtc's contents cannot be zero")
	}

	utc := CleanStringForPlatform(createdUtc.Location().String())
	if strings.Compare(utc, "UTC") != 0 {
		return errors.New("the createdUtc value must be in UTC")
	}

	i.createdUtc = createdUtc

	return nil
}
