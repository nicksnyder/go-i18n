package csv

import (
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"io"
	"csv"
	"os"
)

type Reader struct {
	fieldsep int
}

func NewReader(fieldsep int) *Reader {
	return &Reader{fieldsep: fieldsep}
}

func (cr *Reader) ReadMessages(r io.Reader) ([]msg.Message, os.Error) {
	c := csv.NewReader(r)
	c.Comma = cr.fieldsep
	c.FieldsPerRecord = 4
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
