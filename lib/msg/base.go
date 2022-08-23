package msg

import (
	"github.com/google/uuid"
	"time"
)

type Message interface {
	// FIXME: Fill in interface for messages - BasicMessage (and thus subs) should auto implement
}

type BasicMessage struct {
	// FIXME: Fill in basic message info
	GUID      uuid.UUID             `json:"GUID"`
	From      SourceSystem          `json:"from"`      // What sent the message
	To        []TargetSystem        `json:"to"`        // If len > 0, then only these system should look at the message
	Trace     uuid.UUID             `json:"trace"`     // Used for request tracing
	Generated time.Time             `json:"generated"` // When the message was generated at
	History   map[uuid.UUID]Message `json:"history"`   // Any messages this was in reply to or related - FIXME: Bad idea?
	Key       map[string]string     `json:"key"`       // Key:value pairs used to for routing and should-process decisions
}

func (m BasicMessage) TODO() {

}
