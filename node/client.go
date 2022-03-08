package node

import (
	"net/http"
	"time"
)

type client struct {
	*http.Client
}

func newClient() *client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &client{
		&http.Client{
			Timeout:   time.Minute,
			Transport: t,
		},
	}
}
