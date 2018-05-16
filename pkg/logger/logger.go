package log

import (
	event "github.com/treehouse/simple_event_pusher/pkg/event"
	//"bytes"
	"encoding/json"
	"log"
	"os"
)

type PushLog struct {
	logger *log.Logger
	channel string
	file *os.File
}

func NewPushLog(channel string) *PushLog {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "INFO: ", log.Lshortfile)
	return &PushLog{
		logger: logger,
		channel: channel,
		file: file,
	}
}

func (l PushLog) log(action string, e event.Message) {
	
	record := map[string]string{
		"context": "event_pusher",
		"action": action,
		"channel": e.GetChannel(),
		"event": e.Event(),
		"data": e.Data(),
	}
	
	jData, jErr := json.Marshal(record)
	if jErr != nil {
		log.Println(jErr)
	}
	// bytesLength := bytes.Index(jData, []byte{0})
	jString := string(jData[:]) //bytesLength
	l.logger.Output(2, jString)
}

func (l PushLog) Disconnection() {
	l.log("disconnection", &event.Event{
		ID: "",
		EVENT: "",
		DATA: "",
		CHANNEL: l.channel,
	})
}

func (l PushLog) Connection() {
	l.log("connection", &event.Event{
		ID: "",
		EVENT: "",
		DATA: "",
		CHANNEL: l.channel,
	})
}

func (l PushLog) MessageSent(e event.Message) {
	l.log("message_sent", e)
}

func (l PushLog) Close() {
	if err := l.file.Close(); err != nil {
		log.Fatal(err)
	}
}