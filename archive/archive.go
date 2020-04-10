package archive

import (
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"path/filepath"
	"strings"
)

var (
	ErrBounds = errors.New("Image index out of bounds.")
)

type Archive interface {
	Load(i int, autorotate bool) (*gdk.Pixbuf, error)
	Name(i int) (string, error)
	Len() int
	Close() error
}

const (
	MaxArchiveEntries = 4096 * 64
)

func NewArchive(path string) (Archive, error) {

	switch strings.ToLower(filepath.Ext(path)) {
	case ".zip", ".cbz":
		return NewZip(path)
	case ".7z", ".rar", ".tar", ".tgz", ".tbz2", ".cb7", ".cbr", ".cbt", ".lha":
		// TODO
	case ".gz":
		if strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
			// TODO
		}
	}

	return nil, errors.New("Unknown or unsupported archive type")
}
