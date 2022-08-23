package streams

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/fragforce/fragcenter/lib/msg"
	"github.com/sirupsen/logrus"
)

type IsTrackedFunc func(message msg.Message) bool
type ReactionFunc func(message msg.Message) (*msg.Message, error)

type ActionReaction struct {
	logs.BLog
	action   []IsTrackedFunc // All must validate to True
	reaction ReactionFunc    // then run this
}

func NewActionReaction(log *logrus.Entry, reaction ReactionFunc, actions ...IsTrackedFunc) *ActionReaction {
	return &ActionReaction{
		action:   actions,
		reaction: reaction,
		BLog:     *logs.NewBLog(log),
	}
}
