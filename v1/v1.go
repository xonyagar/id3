package v1

import (
	"errors"
	"fmt"
	"io"
	"github.com/xonyagar/id3/lib"
)

// TagSize is size of ID3v1 and ID3v1.1 tag
const TagSize = 128

// V1 is ID3v1 tag reader
type V1 []byte

// New will read file and return id3v1 tag reader
func New(f io.ReadSeeker) (*V1, error) {
	_, err := f.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	tag := make(V1, TagSize)
	n, err := f.Read(tag)
	if err != nil {
		return nil, err
	}

	if n != TagSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", TagSize, n)
	}

	if string(tag[:3]) != "TAG" {
		return nil, errors.New("no id3v1 tag at the end of file")
	}

	return &tag, nil
}

// Title will return id3v1 title
func (tag V1) Title() string {
	return lib.Trim(string(tag[3:33]))
}

// Artist will return id3v1 artist
func (tag V1) Artist() string {
	return lib.Trim(string(tag[33:63]))
}

// Album will return id3v1 album
func (tag V1) Album() string {
	return lib.Trim(string(tag[63:93]))
}

// Year will return id3v1 year
func (tag V1) Year() string {
	return lib.Trim(string(tag[93:97]))
}

// Comment will return id3v1 or id3v1.1 comment
func (tag V1) Comment() string {
	if tag[125] != byte(0) {
		// V1
		return lib.Trim(string(tag[97:127]))
	} else {
		// V1.1
		return lib.Trim(string(tag[97:125]))
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

	if genre < len(Genres) {
		return Genres[genre]
	}

	return ""
}
