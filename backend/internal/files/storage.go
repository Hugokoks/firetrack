package files

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	root string
}

func NewStorage(root string) *Storage {
	return &Storage{root: root}
}

func (s *Storage) BuildRelativePath(jobID, storedName string) string {
	return filepath.Join("jobs", jobID, storedName)
}

func (s *Storage) BuildAbsolutePath(relPath string) string {
	return filepath.Join(s.root, relPath)
}

func (s *Storage) EnsureDirFor(relPath string) error {
	dir := filepath.Dir(s.BuildAbsolutePath(relPath))
	return os.MkdirAll(dir, 0755)
}

func (s *Storage) Save(relPath string, src io.Reader) (int64, error) {
	absPath := s.BuildAbsolutePath(relPath)

	dst, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		return 0, err
	}

	return written, nil
}

func (s *Storage) Remove(relPath string) error {
	absPath := s.BuildAbsolutePath(relPath)
	if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
