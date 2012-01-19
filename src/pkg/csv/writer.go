package csv

import (
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"io"
	"csv"
	"os"
)

type Writer struct {
	fieldsep int
	crlf     bool
}

func NewWriter(fieldsep int, crlf bool) msg.Writer {
	return &Writer{fieldsep: fieldsep, crlf: crlf}
}

func (w *Writer) WriteMessages(wr io.Writer, msgs []msg.Message) os.Error {
	c := csv.NewWriter(wr)
	defer c.Flush()
	c.Comma = w.fieldsep
	c.UseCRLF = w.crlf
	c.Write(header)
	for _, m := range msgs {
		err := c.Write([]string{m.Id, m.Context, m.Content, m.Translation})
		if err != nil {
			return err
		}
	}
	return nil
}
