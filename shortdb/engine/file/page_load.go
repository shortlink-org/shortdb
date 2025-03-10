package file

import (
	"fmt"
	"os"
	"path/filepath"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	"google.golang.org/protobuf/proto"
)

// PageLoadError represents an error that occurred while loading a page.
type PageLoadError struct {
	Path string
	Err  error
}

func (e *PageLoadError) Error() string {
	return fmt.Sprintf("failed to load page from %s: %v", e.Path, e.Err)
}

func (*File) loadPage(path string) (*page.Page, error) {
	p := page.Page{}

	payload, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, &PageLoadError{Path: path, Err: err}
	}

	if len(payload) != 0 {
		err = proto.Unmarshal(payload, &p)
		if err != nil {
			return nil, &PageLoadError{Path: path, Err: err}
		}
	}

	return &p, nil
}
