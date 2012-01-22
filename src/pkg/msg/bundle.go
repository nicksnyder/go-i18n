package msg

// Bundle is a groups of messages.
type Bundle struct {
	index    map[string]int
	messages []Message
}

// NewBundle creates a new Bundle of Messages.
func NewBundle() *Bundle {
	return &Bundle{index: make(map[string]int, 0), messages: make([]Message, 0)}
}

// AddMessage adds a message to the Bundle.
// If the new message has a non-empty translation, then it will
// replace any message in the Bundle that has the same id.
func (b *Bundle) AddMessage(m Message) {
	m.Id = Id(m.Context, m.Content)
	if i, found := b.index[m.Id]; found && m.Translation != "" {
		b.messages[i] = m
	} else {
		b.messages = append(b.messages, m)
	}
}

// AddMessages adds one or more messages to the Bundle.
// New messages that have a non-empty translation will replace
// messages in the Bundle that have the same id.
func (b *Bundle) AddMessages(msgs []Message) {
	for i, _ := range msgs {
		b.AddMessage(msgs[i])
	}
}

// Messages returns all messages in the Bundle.
func (b *Bundle) Messages() []Message {
	return b.messages
}

// Messages returns all translated messages in the Bundle.
func (b *Bundle) TranslatedMessages() []Message {
	msgs := make([]Message, 0, len(b.messages))
	for _, m := range b.messages {
		if m.Translation != "" {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

// Messages returns all untranslated messages in the Bundle.
func (b *Bundle) UntranslatedMessages() []Message {
	msgs := make([]Message, 0, len(b.messages))
	for _, m := range b.messages {
		if m.Translation == "" {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

// Len implementes sort.Interface
func (b *Bundle) Len() int {
	return len(b.messages)
}

// Less implementes sort.Interface
func (b *Bundle) Less(i, j int) bool {
	mi := b.messages[i]
	mj := b.messages[j]
	if mi.Content < mj.Content {
		return true
	}
	if mi.Content > mj.Content {
		return false
	}
	return mi.Context < mj.Context
}

// Swap implementes sort.Interface
func (b *Bundle) Swap(i, j int) {
	b.messages[i], b.messages[j] = b.messages[j], b.messages[i]
}
