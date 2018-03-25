package id3

import (
	"errors"
	"fmt"
	"io"
)

// V1TagSize is size of ID3v1 and ID3v1.1 tag
const V1TagSize = 128

// V1 is ID3v1 tag reader
type V1 []byte

// NewID3V1 will read file and return id3v1 tag reader
func NewID3V1(f io.ReadSeeker) (*V1, error) {
	_, err := f.Seek(-V1TagSize, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	tag := make(V1, V1TagSize)
	n, err := f.Read(tag)
	if err != nil {
		return nil, err
	}

	if n != V1TagSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", V1TagSize, n)
	}

	if string(tag[:3]) != "TAG" {
		return nil, errors.New("no id3v1 tag at the end of file")
	}

	return &tag, nil
}

// Title will return id3v1 title
func (tag V1) Title() string {
	return trim(string(tag[3:33]))
}

// Artist will return id3v1 artist
func (tag V1) Artist() string {
	return trim(string(tag[33:63]))
}

// Album will return id3v1 album
func (tag V1) Album() string {
	return trim(string(tag[63:93]))
}

// Year will return id3v1 year
func (tag V1) Year() string {
	return trim(string(tag[93:97]))
}

// Comment will return id3v1 or id3v1.1 comment
func (tag V1) Comment() string {
	if tag[125] != byte(0) {
		// V1
		return trim(string(tag[97:127]))
	} else {
		// V1.1
		return trim(string(tag[97:125]))
	}
}

// AlbumTrack will return id3v1.1 album track
func (tag V1) AlbumTrack() string {
	if tag[125] == byte(0) {
		return fmt.Sprintf("%d", int(tag[126]))
	} else {
		return ""
	}
}

// Genre will return id3v1 genre title
func (tag V1) Genre() string {

	genre := int(tag[127])

	if genre < len(V1Genres) {
		return V1Genres[genre]
	}

	return ""
}
