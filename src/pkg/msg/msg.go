package msg

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// Reader is the interface that provides a method
// to read messages from an io.ReadSeeker.
type Reader interface {
	ReadMessages(rs io.ReadSeeker) ([]Message, os.Error)
}

// Writer is the interface that provides a method
// to write messages to an io.Writer.
type Writer interface {
	WriteMessages(w io.Writer, msgs []Message) os.Error
}

// Message represents a serializable message
// for the purpose of reading from and writing to files.
type Message struct {
	Id          string
	Context     string
	Content     string
	Translation string
}

func Id(context, content string) string {
	h := md5.New()
	h.Write([]byte(context))
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum())
}
