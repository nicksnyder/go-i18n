package message

// Message represents a serializable message
// for the purpose of reading from and writing to files.
type Message struct {
	Context     string
	Content     string
	Translation string
}
