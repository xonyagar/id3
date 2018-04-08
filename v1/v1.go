package v1

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// TagSize is size of ID3v1 and ID3v1.1 tag
const TagSize = 128

const Version10 = "ID3v1"

const Version11 = "ID3v1.1"

const TagText = "TAG"

// Tag is ID3v1 tag reader
type Tag struct {
	version    string
	title      string
	artist     string
	album      string
	year       string
	comment    string
	albumTrack int
	genreIndex int
}

var ErrTagNotFound = errors.New("no id3v1 tag at the end of file")

// Remove id3v1 tag from end of file
func Remove(f *os.File) error {
	_, err := f.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return err
	}

	b := make([]byte, TagSize)
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	if n != TagSize {
		return fmt.Errorf("must read '%d' bytes, but read '%d'", TagSize, n)
	}

	if string(b[:3]) != TagText {
		return ErrTagNotFound
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	return f.Truncate(fi.Size() - TagSize)
}

// New will read file and return id3v1 tag reader
func New(f io.ReadSeeker) (*Tag, error) {
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

	tag := Tag{}

	if string(b[:3]) != TagText {
		tag.Clear()
		return &tag, ErrTagNotFound
	}

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

	if b[125] == 0 && b[126] != 0 {
		// V1.1
		tag.version = Version11
		for i := 97; i < 125; i++ {
			if b[i] == 0 {
				break
			}
			tag.comment += string(b[i])
		}

		tag.albumTrack = int(b[126])
	} else {
		// V1
		tag.version = Version10
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
func (tag Tag) Title() string {
	return tag.title
}

// SetTitle will set id3v1 title
func (tag *Tag) SetTitle(title string) {
	tag.title = title
}

// Artist will return id3v1 artist
func (tag Tag) Artist() string {
	return tag.artist
}

// SetArtist will set id3v1 artist
func (tag *Tag) SetArtist(artist string) {
	tag.artist = artist
}

// Album will return id3v1 album
func (tag Tag) Album() string {
	return tag.album
}

// SetAlbum will set id3v1 album
func (tag *Tag) SetAlbum(album string) {
	tag.album = album
}

// Year will return id3v1 year
func (tag Tag) Year() string {
	return tag.year
}

// SetYear will set id3v1 year
func (tag *Tag) SetYear(year string) {
	tag.year = year
}

// Comment will return id3v1 or id3v1.1 comment
func (tag Tag) Comment() string {
	return tag.comment
}

// SetComment will set id3v1 comment
func (tag *Tag) SetComment(comment string) {
	tag.comment = comment
}

// AlbumTrack will return id3v1.1 album track
func (tag Tag) AlbumTrack() int {
	return tag.albumTrack
}

// SetAlbumTrack will set id3v1 album track
func (tag *Tag) SetAlbumTrack(albumTrack int) {
	tag.albumTrack = albumTrack
	tag.SetVersionTo11()
}

// GenreIndex will return id3v1 genre index
func (tag Tag) GenreIndex() int {
	return tag.genreIndex
}

// SetGenreIndex will set id3v1 genre index
func (tag *Tag) SetGenreIndex(genreIndex int) {
	if genreIndex < len(Genres) {
		tag.genreIndex = genreIndex
	} else {
		tag.genreIndex = 255
	}
}

// Genre will return id3v1 genre title
func (tag Tag) Genre() string {
	if tag.genreIndex < len(Genres) {
		return Genres[tag.genreIndex]
	}

	return "Unknown"
}

// SetGenre will set id3v1 genre by title
func (tag *Tag) SetGenre(genre string) {
	for i := range Genres {
		if Genres[i] == genre {
			tag.genreIndex = i
			return
		}
	}

	tag.genreIndex = 255
}

// Version will return id3 version
func (tag Tag) Version() string {
	return tag.version
}

// SetVersionTo10 will set id3 version to v1
func (tag *Tag) SetVersionTo10() {
	tag.version = Version10
}

// SetVersionTo11 will set id3 version to v1.1
func (tag *Tag) SetVersionTo11() {
	tag.version = Version11
}

// SetVersionTo11 will set id3 version to v1.1
func (tag *Tag) Clear() {
	tag.version = Version10
	tag.title = ""
	tag.artist = ""
	tag.album = ""
	tag.year = ""
	tag.comment = ""
	tag.albumTrack = 0
	tag.genreIndex = 255
}

func (tag Tag) Write(f io.ReadWriteSeeker) error {
	_, err := f.Seek(-TagSize, io.SeekEnd)
	if err != nil {
		return err
	}

	b := make([]byte, TagSize)
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	if n != TagSize {
		return fmt.Errorf("must read '%d' bytes, but read '%d'", TagSize, n)
	}

	if string(b[:3]) == TagText {
		_, err := f.Seek(-TagSize, io.SeekEnd)
		if err != nil {
			return err
		}
	}

	// Create Tag Binary
	tb := make([]byte, TagSize)

	// Tag Text
	for i := 0; i < 3 && i < len(TagText); i++ {
		tb[i] = TagText[i]
	}

	// Title
	for i := 3; i < 33 && i-3 < len(tag.title); i++ {
		tb[i] = tag.title[i-3]
	}

	// Artist
	for i := 33; i < 63 && i-33 < len(tag.artist); i++ {
		tb[i] = tag.artist[i-33]
	}

	// Album
	for i := 63; i < 93 && i-63 < len(tag.album); i++ {
		tb[i] = tag.album[i-63]
	}

	// Year
	for i := 93; i < 97 && i-93 < len(tag.year); i++ {
		tb[i] = tag.year[i-93]
	}

	if tag.version == Version10 {
		// Comment
		for i := 97; i < 127 && i-97 < len(tag.comment); i++ {
			tb[i] = tag.comment[i-97]
		}
	}

	if tag.version == Version11 {
		// Comment
		for i := 97; i < 125 && i-97 < len(tag.comment); i++ {
			tb[i] = tag.comment[i-97]
		}

		// Album Track
		tb[125] = 0
		tb[126] = byte(tag.albumTrack)
	}

	// Genre
	tb[127] = byte(tag.genreIndex)

	_, err = f.Write(tb)
	return err
}
