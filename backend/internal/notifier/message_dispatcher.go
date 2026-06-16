package notifier

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, msg *Message) error

type MessageDispatcher struct {
	handlers map[string]Handler
}

func NewMessageDispatcher() *MessageDispatcher {
	return &MessageDispatcher{
		handlers: make(map[string]Handler),
	}
}

func (d *MessageDispatcher) Register(mtype MessageType, handler Handler) {
	d.handlers[string(mtype)] = handler
}

func (d *MessageDispatcher) Dispatch(ctx context.Context, msg *Message) error {
	handler, ok := d.handlers[string(msg.Type)]
	if !ok {
		return fmt.Errorf("no handler for message :%v", msg.Type)
	}

	return handler(ctx, msg)
}
