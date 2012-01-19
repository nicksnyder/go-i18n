package csv

import (
	"bufio"
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"io"
	"csv"
	"os"
)

type Reader struct {

}

func NewReader() msg.Reader {
	return &Reader{}
}

func (r *Reader) ReadMessages(rs io.ReadSeeker) ([]msg.Message, os.Error) {
	// Peek ahead to detect the separator character that was used to encode the header
	_, err := rs.Seek(0, len(header[0]))
	if err != nil {
		return nil, err
	}

	// The next rune should be the separator
	rd, err := bufio.NewReaderSize(rs, 8)
	if err != nil {
		return nil, err
	}
	fieldsep, _, err := rd.ReadRune()

	// Go back to the beginning of the file to parse as csv
	_, err = rs.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	c := csv.NewReader(rs)
	c.Comma = fieldsep
	c.FieldsPerRecord = len(header)
	c.TrailingComma = true

	lines, err := c.ReadAll()
	if err != nil {
		return nil, err
	}

	// Trim header
	lines = lines[1:]

	m := make([]msg.Message, len(lines))
	for i, line := range lines {
		m[i].Id = msg.Id(line[1], line[2])
		m[i].Context = line[1]
		m[i].Content = line[2]
		m[i].Translation = line[3]
	}
	return m, nil
}
