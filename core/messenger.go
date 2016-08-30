package core

import (
	"github.com/docker/engine-api/types/events"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/golang/protobuf/proto"
)

// Messenger singleton
type Messenger struct {
}

var (
	nats Nats
)

// Init the Messenger
func (Messenger *Messenger) Init() error {
	return nats.Connect()
}

func (Messenger *Messenger) sendLog(entry logs.LogEntry) error {
	encoded, err := proto.Marshal(&entry)
	if err != nil {
		return err
	}
	_, err = nats.Publish("amp-logs", encoded)
	return err
}

func (Messenger *Messenger) sendEvent(event events.Message) error {
	//encoded, err := proto.Marshal(&event)
	//if err != nil {
	//	return err
	//}
	//_, err = nats.Publish("amp-events", encoded)
	return nil
}

// Close the Messenger
func (Messenger *Messenger) Close() {
	// TODO
}
