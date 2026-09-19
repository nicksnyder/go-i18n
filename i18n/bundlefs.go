package i18n

import (
	"io/fs"

	"golang.org/x/text/language"
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

// LoadMessageFileFSWithTag is like LoadMessageFileFS but uses tag as the file
// language instead of inferring it from path.
func (b *Bundle) LoadMessageFileFSWithTag(fsys fs.FS, path string, tag language.Tag) (*MessageFile, error) {
	buf, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}
	return b.ParseMessageFileBytesWithTag(buf, path, tag)
}
