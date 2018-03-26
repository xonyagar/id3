package v1

import (
	"errors"
	"fmt"
	"io"
)

// TagSize is size of ID3v1 and ID3v1.1 tag
const TagSize = 128

// V1 is ID3v1 tag reader
type V1 struct {
	title      string
	artist     string
	album      string
	year       string
	comment    string
	albumTrack string
	genreIndex int
}

var ErrTagNotFound = errors.New("no id3v1 tag at the end of file")

// New will read file and return id3v1 tag reader
func New(f io.ReadSeeker) (*V1, error) {
	_, err := f.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	b := make([]byte, TagSize)
	n, err := f.Read(b)
	if err != nil {
		return nil, err
	}

	if n != TagSize {
		return nil, fmt.Errorf("must read '%d' bytes, but read '%d'", TagSize, n)
	}

	if string(b[:3]) != "TAG" {
		return nil, ErrTagNotFound
	}

	tag := V1{}
	for i := 3; i < 33; i++ {
		if b[i] == 0 {
			break
		}
		tag.title += string(b[i])
	}
	for i := 33; i < 63; i++ {
		if b[i] == 0 {
			break
		}
		tag.artist += string(b[i])
	}
	for i := 63; i < 93; i++ {
		if b[i] == 0 {
			break
		}
		tag.album += string(b[i])
	}
	for i := 93; i < 97; i++ {
		if b[i] == 0 {
			break
		}
		tag.year += string(b[i])
	}

	if b[125] == 0 {
		// V1.1
		for i := 97; i < 125; i++ {
			if b[i] == 0 {
				break
			}
			tag.comment += string(b[i])
		}

		tag.albumTrack = fmt.Sprintf("%d", int(b[126]))
	} else {
		// V1
		for i := 97; i < 127; i++ {
			if b[i] == 0 {
				break
			}
			tag.comment += string(b[i])
		}
	}

	tag.genreIndex = int(b[127])

	return &tag, nil
}

// Title will return id3v1 title
func (tag V1) Title() string {
	return tag.title
}

// Artist will return id3v1 artist
func (tag V1) Artist() string {
	return tag.artist
}

// Album will return id3v1 album
func (tag V1) Album() string {
	return tag.album
}

// Year will return id3v1 year
func (tag V1) Year() string {
	return tag.year
}

// Comment will return id3v1 or id3v1.1 comment
func (tag V1) Comment() string {
	return tag.comment
}

// AlbumTrack will return id3v1.1 album track
func (tag V1) AlbumTrack() string {
	return tag.albumTrack
}

// Genre will return id3v1 genre title
func (tag V1) Genre() string {
	if tag.genreIndex < len(Genres) {
		return Genres[tag.genreIndex]
	}

	return ""
}
