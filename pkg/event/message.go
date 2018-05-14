package event

// Boilerplate for eventsource package
type Message struct {
	EventStr string `json:"event"`
	Channel  string `json:"channel"`
	DataStr  string `json:"data"`
}

func (m *Message) Id() string    { return "" }
func (m *Message) Event() string { return m.EventStr }
func (m *Message) Data() string  { return m.DataStr }
