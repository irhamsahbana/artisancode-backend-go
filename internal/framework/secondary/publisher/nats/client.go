package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type client struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func newClient(url string) (*client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &client{
		nc: nc,
		js: js,
	}, nil
}

func (c *client) Close() error {
	if c.nc != nil {
		c.nc.Close()
	}

	return nil
}
