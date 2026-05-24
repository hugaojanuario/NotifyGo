package channel

import "context"

type Event struct {
	Topic   string
	Payload string
}

type Channel interface {
	Send(ctx context.Context, event Event, cfg ChannelConfig) error
}

type Registry struct {
	channels map[ChannelType]Channel
}

func NewRegistry() *Registry {
	return &Registry{channels: make(map[ChannelType]Channel)}
}

func (r *Registry) Register(t ChannelType, c Channel) {
	r.channels[t] = c
}

func (r *Registry) Get(t ChannelType) (Channel, bool) {
	c, ok := r.channels[t]
	return c, ok
}
