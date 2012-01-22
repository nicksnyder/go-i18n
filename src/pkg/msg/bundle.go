package msg

// Bundle is a groups of messages.
type Bundle struct {
	filter   Filter
	index    map[string]int
	messages []Message
}

type Filter uint

const (
	AllMessages Filter = iota
	TranslatedOnly
	UntranslatedOnly
)

// NewBundle creates a new Bundle of Messages.
func NewBundle() *Bundle {
	return &Bundle{index: make(map[string]int, 0), messages: make([]Message, 0)}
}

// AddMessage adds a new Message to the Bundle.
// If the new Message has a non-empty translation, then it will overwrite any
// existing Message in the Bundle that has the same id.
func (b *Bundle) AddMessage(m Message) {
	if b.filter == TranslatedOnly && m.Translation == "" {
		return
	}
	m.Id = Id(m.Context, m.Content)
	if i, found := b.index[m.Id]; found && m.Translation != "" {
		b.messages[i] = m
	} else {
		b.messages = append(b.messages, m)
	}
}

// Messages returns the messages that match the filter.
func (b *Bundle) Messages(filter Filter) []Message {
	if filter == AllMessages {
		return b.messages
	}

	msgs := make([]Message, 0, len(b.messages))
	for _, m := range b.messages {
		if filter == UntranslatedOnly && m.Translation == "" ||
			filter == TranslatedOnly && m.Translation != "" {
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
