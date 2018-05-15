package push

import (
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	"net/http"
)

type Connection interface {
	ServePUSH(http.ResponseWriter, *http.Request)
	Send(event.Message)
	Channel() string
	Close()
	Msgs()
}
