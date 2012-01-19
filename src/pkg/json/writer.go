package json

import (
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"io"
	"json"
	"os"
)

type Writer struct {

}

func NewWriter() msg.Writer {
	return &Writer{}
}

func (w *Writer) WriteMessages(iw io.Writer, m []msg.Message) os.Error {
	json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = iw.Write(json)
	return err
}
