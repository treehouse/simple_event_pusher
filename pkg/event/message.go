package event

type Message interface {
	Id() string
	Event() string
	Data() string
	GetChannel() string
}

type Event struct {
	ID      string
	EVENT   string
	DATA    string
	CHANNEL string
}

func (e Event) GetChannel() string { return e.CHANNEL }
func (e Event) Id() string         { return e.ID }
func (e Event) Event() string      { return e.EVENT }
func (e Event) Data() string       { return e.DATA }

// Boilerplate for eventsource package
// type JSONMessage struct {
// 	EventStr string `json:"event"`
// 	Channel  string `json:"channel"`
// 	DataStr  string `json:"data"`
// }

// func (m *JSONMessage) Id() string    { return "" }
// func (m *JSONMessage) Event() string { return m.EventStr }
// func (m *JSONMessage) Data() string  { return m.DataStr }
// func (m *JSONMessage) GetChannel() string  { return m.Channel }
