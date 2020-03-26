package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

// Message represents the information an event contains
type DockerEvent struct {
	ID     string `json:"id,omitempty"`
	From   string `json:"from,omitempty"`
	Type   string
	Action string
	Actor  struct {
		ID         string
		Attributes map[string]string
	}
	Scope string `json:"scope,omitempty"`
	TimeNano int64 `json:"timeNano,omitempty"`
}

func eventsReader(ctx context.Context, yamlConfig YamlConfig)  (<-chan DockerEvent, <-chan error)  {
	messages := make(chan DockerEvent)
	errs := make(chan error, 1)

	started := make(chan struct{})

	go func() {
		defer close(errs)

		httpc := http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", "/var/run/docker.sock")
				},
			},
		}

		response, err := httpc.Get("http://docker/v1.40/events")

		if err != nil {
			close(started)
			errs <- err
			return
		}

		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)

		close(started)
		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				var event DockerEvent
				if err := decoder.Decode(&event); err != nil {
					errs <- err
					return
				}

				select {
				case messages <- event:
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				}
			}
		}
	}()
	<-started

	return messages, errs
}