package id3

import (
	"errors"
	"fmt"
	"image"
	"io"
	"strconv"

	v1 "github.com/xonyagar/id3/v1"
	v22 "github.com/xonyagar/id3/v22"
	v23 "github.com/xonyagar/id3/v23"
	v24 "github.com/xonyagar/id3/v24"
)

type ID3 struct {
	V1  *v1.Tag
	V22 *v22.Tag
	V23 *v23.Tag
	V24 *v24.Tag
}

func New(f io.ReadSeeker) (*ID3, error) {
	tag := new(ID3)

	var err error

	tag.V1, err = v1.New(f)
	if err != nil && !errors.Is(err, v1.ErrTagNotFound) {
		return nil, fmt.Errorf("error on new v1: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error on seek: %w", err)
	}

	tag.V22, err = v22.New(f)
	if err != nil && !errors.Is(err, v22.ErrTagNotFound) {
		return nil, fmt.Errorf("error on new v2.2: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error on seek: %w", err)
	}

	tag.V23, err = v23.New(f)
	if err != nil && !errors.Is(err, v23.ErrTagNotFound) {
		return nil, fmt.Errorf("error on new v2.3: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error on seek: %w", err)
	}

	tag.V24, err = v24.New(f)
	if err != nil && !errors.Is(err, v24.ErrTagNotFound) {
		return nil, fmt.Errorf("error on new v2.4: %w", err)
	}

	return tag, nil
}

func (t ID3) Title() string {
	if t.V24 != nil {
		if title := t.V24.Title(); title != "" {
			return title
		}
	}

	if t.V23 != nil {
		if title := t.V23.Title(); title != "" {
			return title
		}
	}

	if t.V22 != nil {
		if title := t.V22.Title(); title != "" {
			return title
		}
	}

	if t.V1 != nil {
		if title := t.V1.Title(); title != "" {
			return title
		}
	}

	return ""
}

func (t ID3) Album() string {
	if t.V24 != nil {
		if album := t.V24.Album(); album != "" {
			return album
		}
	}

	if t.V23 != nil {
		if album := t.V23.Album(); album != "" {
			return album
		}
	}

	if t.V22 != nil {
		if album := t.V22.Album(); album != "" {
			return album
		}
	}

	if t.V1 != nil {
		if album := t.V1.Album(); album != "" {
			return album
		}
	}

	return ""
}

func (t ID3) AlbumArtists() []string {
	if t.V24 != nil {
		if albumArtists := t.V24.AlbumArtists(); len(albumArtists) > 0 {
			return albumArtists
		}
	}

	if t.V23 != nil {
		if albumArtists := t.V23.AlbumArtists(); len(albumArtists) > 0 {
			return albumArtists
		}
	}

	if t.V22 != nil {
		if albumArtists := t.V22.AlbumArtists(); len(albumArtists) > 0 {
			return albumArtists
		}
	}

	if t.V1 != nil {
		if artist := t.V1.Artist(); artist != "" {
			return []string{artist}
		}
	}

	return []string{}
}

func (t ID3) Artists() []string {
	if t.V24 != nil {
		if artists := t.V24.Artists(); len(artists) > 0 {
			return artists
		}
	}

	if t.V23 != nil {
		if artists := t.V23.Artists(); len(artists) > 0 {
			return artists
		}
	}

	if t.V22 != nil {
		if artists := t.V22.Artists(); len(artists) > 0 {
			return artists
		}
	}

	if t.V1 != nil {
		if artist := t.V1.Artist(); artist != "" {
			return []string{artist}
		}
	}

	return []string{}
}

func (t ID3) TrackNumberAndPosition() (int, int) {
	if t.V24 != nil {
		if a, b := t.V24.TrackNumberAndPosition(); a != 0 {
			return a, b
		}
	}

	if t.V23 != nil {
		if a, b := t.V23.TrackNumberAndPosition(); a != 0 {
			return a, b
		}
	}

	if t.V22 != nil {
		if a, b := t.V22.TrackNumberAndPosition(); a != 0 {
			return a, b
		}
	}

	if t.V1 != nil {
		if s := t.V1.AlbumTrack(); s != "" {
			a, err := strconv.Atoi(s)
			if err != nil {
				return a, 0
			}
		}
	}

	return 0, 0
}

func (t ID3) Year() string {
	if t.V24 != nil {
		if year := t.V24.Year(); year != "" {
			return year
		}
	}

	if t.V23 != nil {
		if year := t.V23.Year(); year != "" {
			return year
		}
	}

	if t.V22 != nil {
		if year := t.V22.Year(); year != "" {
			return year
		}
	}

	if t.V1 != nil {
		if year := t.V1.Year(); year != "" {
			return year
		}
	}

	return ""
}

type AttachedPicture interface {
	Image() (image.Image, error)
}

func (t ID3) AttachedPictures() []AttachedPicture {
	if t.V24 != nil {
		if pics := t.V24.AttachedPictures(); len(pics) > 0 {
			res := make([]AttachedPicture, len(pics))
			for i := range pics {
				res[i] = pics[i]
			}

			return res
		}
	}

	if t.V23 != nil {
		if pics := t.V23.AttachedPictures(); len(pics) > 0 {
			res := make([]AttachedPicture, len(pics))
			for i := range pics {
				res[i] = pics[i]
			}

			return res
		}
	}

	if t.V22 != nil {
		if pics := t.V22.AttachedPictures(); len(pics) > 0 {
			res := make([]AttachedPicture, len(pics))
			for i := range pics {
				res[i] = pics[i]
			}

			return res
		}
	}

	return []AttachedPicture{}
}

func (t ID3) Genres() []string {
	if t.V24 != nil {
		if genres := t.V24.Genres(); len(genres) > 0 {
			return genres
		}
	}

	if t.V23 != nil {
		if genres := t.V23.Genres(); len(genres) > 0 {
			return genres
		}
	}

	if t.V22 != nil {
		if genres := t.V22.Genres(); len(genres) > 0 {
			return genres
		}
	}

	if t.V1 != nil {
		if artist := t.V1.Genre(); artist != "" {
			return []string{artist}
		}
	}

	return []string{}
}
