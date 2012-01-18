package msg

import (
	"crypto/md5"
	"encoding/hex"
)

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
