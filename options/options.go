package options

import (
	"path/filepath"
)

type Options struct {
	Prefix  FilePath
	PubRing FilePath
	SecRing FilePath
}

type FilePath string

func (fp *FilePath) String() string {
	return string(*fp)
}

func (fp *FilePath) Set(val string) (err error) {
	fps, err := filepath.Abs(val)
	*fp = FilePath(fps)
	return err
}
