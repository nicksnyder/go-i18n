package json

import (
	"github.com/nicksnyder/go-i18n/src/pkg/message"
	"io"
	"io/ioutil"
	"json"
	"os"
)

type Formatter struct {

}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) ReadMessages(r io.Reader) ([]message.Message, os.Error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	m := make([]message.Message, 0, 100)
	err = json.Unmarshal(data, m)
	return m, err
}

func (f *Formatter) WriteMessages(w io.Writer, m []message.Message) os.Error {
	json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(json)
	return err
}
