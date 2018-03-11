package id3

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// V11 is ID3v1.1 tag reader
type V11 []byte

// NewID3V11 will read file and return id3v1.1 tag reader
func NewID3V11(f io.ReadSeeker) (*V11, error) {
	_, err := f.Seek(-V1TagSize, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	tag := make(V11, V1TagSize)
	n, err := f.Read(tag)
	if err != nil {
		return nil, err
	}

	if n != V1TagSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", V1TagSize, n)
	}

	if string(tag[:3]) != "TAG" {
		return nil, errors.New("no id3v1.1 tag at the end of file")
	}

	return &tag, nil
}

// Title will return id3v1.1 title
func (tag V11) Title() string {
	return strings.Trim(string(tag[3:33]), string(uint(0)))
}

// Artist will return id3v1.1 artist
func (tag V11) Artist() string {
	return strings.Trim(string(tag[33:63]), string(uint(0)))
}

// Album will return id3v1.1 album
func (tag V11) Album() string {
	return strings.Trim(string(tag[63:93]), string(uint(0)))
}

// Year will return id3v1.1 year
func (tag V11) Year() string {
	return strings.Trim(string(tag[93:97]), string(uint(0)))
}

// Comment will return id3v1.1 comment
func (tag V11) Comment() string {
	return strings.Trim(string(tag[97:125]), string(uint(0)))
}

// AlbumTrack will return id3v1.1 album track
func (tag V11) AlbumTrack() string {
	return fmt.Sprintf("%d", int(tag[126]))
}

// Genre will return id3v1.1 genre title
func (tag V11) Genre() string {

	genre := int(tag[127])

	if genre < len(V1Genres) {
		return V1Genres[genre]
	}

	return ""
}
