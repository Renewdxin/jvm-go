package classpath

import (
	"os"
	"path/filepath"
)

type DirEntry struct {
	absDir string
}

func newDirEntry(path string) *DirEntry {
    return &DirEntry{path}
}

func (self *DirEntry) readClass(className string) ([]byte, Entry, error) {
	fileName := filepath.Join(self.absDir, className)
	data, err := os.ReadFile(fileName)
	return data, self, err
}

func (self *DirEntry) String() string {
    return self.absDir
}

