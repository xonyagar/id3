package id3

import (
	"io"

	"github.com/xonyagar/id3/v1"
	"github.com/xonyagar/id3/v22"
	"github.com/xonyagar/id3/v23"
	"github.com/xonyagar/id3/v24"
)

type ID3 struct {
	V1  *v1.V1
	V22 *v22.V22
	V23 *v23.V23
	V24 *v24.V24
}

func New(f io.ReadSeeker) (*ID3, error) {
	tag := ID3{}
	var err error

	tag.V1, err = v1.New(f)
	if err != v1.ErrTagNotFound {
		return nil, err
	}

	f.Seek(0, 0)
	tag.V22, err = v22.New(f)
	if err != v22.ErrTagNotFound {
		return nil, err
	}

	f.Seek(0, 0)
	tag.V23, err = v23.New(f)
	if err != v23.ErrTagNotFound {
		return nil, err
	}

	f.Seek(0, 0)
	tag.V24, err = v24.New(f)
	if err != v24.ErrTagNotFound {
		return nil, err
	}

	return &tag, nil
}
