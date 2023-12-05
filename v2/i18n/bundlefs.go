package i18n

import (
	"io/fs"
)

// LoadMessageFileFS is like LoadMessageFile but instead of reading from the
// hosts operating system's file system it reads from the fs file system.
func (b *Bundle) LoadMessageFileFS(fsys fs.FS, path string) (*MessageFile, error) {
	buf, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	return b.ParseMessageFileBytes(buf, path)
}
