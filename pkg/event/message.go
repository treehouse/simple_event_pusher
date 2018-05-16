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