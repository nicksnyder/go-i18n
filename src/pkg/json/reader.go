package json

import (
	"github.com/nicksnyder/go-i18n/src/pkg/msg"
	"io"
	"io/ioutil"
	"json"
	"os"
)

type Reader struct {

}

func NewReader() msg.Reader {
	return &Reader{}
}

func (r *Reader) ReadMessages(ir io.Reader) ([]msg.Message, os.Error) {
	data, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, err
	}
	m := make([]msg.Message, 0, 100)
	err = json.Unmarshal(data, m)
	return m, err
}
