package fs

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/lindell/multi-gitter/internal/template"

	"github.com/pkg/errors"
)

var fileWriterCache map[string]*os.File
var mutex sync.Mutex

// NopWriter is a writer that does nothing
type NopWriter struct{}

func (nw NopWriter) Write(bb []byte) (int, error) {
	return len(bb), nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func FileOutput(path string, std io.Writer) (io.WriteCloser, error) {
	// prevent creation of multiple file streams per path
	mutex.Lock()
	defer mutex.Unlock()

	if path == "-" {
		return nopCloser{std}, nil
	}

	if cached, ok := fileWriterCache[path]; !ok {
		return cached, nil
	}

	// create folders if necessary
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 666)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create folder %s", dir)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %s", path)
	}

	fileWriterCache[path] = file
	return file, nil
}

func FileOutputTemplated(pattern string, repo scm.Repository, std io.Writer) (io.WriteCloser, error) {
	stdoutPath, err := template.Template(pattern, repo)

	output, err := FileOutput(stdoutPath, std)
	if err != nil {
		return nil, err
	}
	return output, nil
}
